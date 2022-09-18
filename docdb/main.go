package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/cockroachdb/pebble"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	db      *pebble.DB
	indexDb *pebble.DB
	port    int
}

func newServer(database string, port int) (*server, error) {
	db, err := pebble.Open(database, &pebble.Options{})
	if err != nil {
		return nil, err
	}

	indexDb, err := pebble.Open(database+".index", &pebble.Options{})
	if err != nil {
		return nil, err
	}

	return &server{db, indexDb, port}, nil
}

func main() {
	s, err := newServer("docdb.data", 8080)
	if err != nil {
		log.Fatal(err)
	}
	defer s.db.Close()

	s.reindex()

	router := httprouter.New()
	router.POST("/docs", s.addDocument)
	router.GET("/docs", s.searchDocuments)
	router.GET("/docs/:id", s.getDocument)

	strPort := strconv.Itoa(s.port)
	log.Println("Listening on " + strPort)
	log.Fatal(http.ListenAndServe(":"+strPort, router))
}

func (s server) addDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	var document map[string]any
	err := dec.Decode(&document)
	if err != nil {
		jsonResponse(w, nil, err)
		return

	}

	// New unique id for the document
	id := uuid.New().String()

	s.index(id, document)

	bs, err := json.Marshal(document)
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	err = s.db.Set([]byte(id), bs, pebble.Sync)
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	jsonResponse(w, map[string]any{
		"id": id,
	}, nil)
}

func (s server) lookup(pathValue string) ([]string, error) {
	idsString, closer, err := s.indexDb.Get([]byte(pathValue))
	if err != nil && err != pebble.ErrNotFound {
		return nil, fmt.Errorf("Could not look up pathvalue [%#v]: %s", pathValue, err)
	}
	if closer != nil {
		defer closer.Close()
	}

	if len(idsString) == 0 {
		return nil, nil
	}
	return strings.Split(string(idsString), ","), nil
}

func (s server) searchDocuments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q, err := parseQuery(r.URL.Query().Get("q"))
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	isRange := false
	idsArgumentCount := map[string]int{}
	nonRangeArguments := 0
	for _, argument := range q.ands {
		if argument.op == "=" {
			nonRangeArguments++

			ids, err := s.lookup(fmt.Sprintf("%s=%v", strings.Join(argument.key, "."), argument.value))
			if err != nil {
				jsonResponse(w, nil, err)
				return
			}

			for _, id := range ids {
				_, ok := idsArgumentCount[id]
				if !ok {
					idsArgumentCount[id] = 0
				}
				idsArgumentCount[id]++
			}
		} else {
			isRange = false
		}
	}

	var idsInAll []string
	for id, count := range idsArgumentCount {
		if count == nonRangeArguments {
			idsInAll = append(idsInAll, id)
		}
	}

	var documents []map[string]any
	if r.URL.Query().Get("skipIndex") == "true" {
		idsInAll = nil
	}

	if len(idsInAll) > 0 {
		for _, id := range idsInAll {
			document, err := s.getDocumentById([]byte(id))
			if err != nil {
				jsonResponse(w, nil, err)
				return
			}

			if !isRange || q.match(document) {
				documents = append(documents, map[string]any{
					"id":   string(id),
					"body": document,
				})
			}
		}
	} else {
		iter := s.db.NewIter(nil)
		defer iter.Close()

		for iter.First(); iter.Valid(); iter.Next() {
			var document map[string]any
			err = json.Unmarshal(iter.Value(), &document)
			if err != nil {
				jsonResponse(w, nil, err)
				return
			}

			if q.match(document) {
				documents = append(documents, map[string]any{
					"id":   string(iter.Key()),
					"body": document,
				})
			}
		}
	}

	jsonResponse(w, map[string]any{
		"documents": documents,
		"count":     len(documents),
	}, nil)
}

func (s server) getDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	document, err := s.getDocumentById([]byte(id))
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	jsonResponse(w, map[string]any{
		"document": document,
	}, nil)
}

func (s server) getDocumentById(id []byte) (map[string]any, error) {
	valBytes, closer, err := s.db.Get(id)
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	var document map[string]any
	err = json.Unmarshal(valBytes, &document)
	return document, nil
}

func (s server) reindex() {
	iter := s.db.NewIter(nil)
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		var document map[string]any
		err := json.Unmarshal(iter.Value(), &document)
		if err != nil {
			log.Printf("Unable to parse bad document, %s: %s", string(iter.Key()), err)
		}
		s.index(string(iter.Key()), document)
	}
}

func (s server) index(id string, document map[string]any) {
	pv := getPathValues(document, "")

	for _, pathValue := range pv {
		idsString, closer, err := s.indexDb.Get([]byte(pathValue))
		if err != nil && err != pebble.ErrNotFound {
			log.Printf("Could not look up pathvalue [%#v]: %s", document, err)
		}

		if len(idsString) == 0 {
			idsString = []byte(id)
		} else {
			ids := strings.Split(string(idsString), ",")

			found := false
			for _, existingId := range ids {
				if id == existingId {
					found = true
				}
			}

			if !found {
				idsString = append(idsString, []byte(","+id)...)
			}
		}

		if closer != nil {
			err := closer.Close()
			if err != nil {
				log.Printf("Could not close: %s", err)
			}
		}

		err = s.indexDb.Set([]byte(pathValue), idsString, pebble.Sync)
		if err != nil {
			log.Printf("Could not update index: %s", err)
		}
	}
}

func getPathValues(obj map[string]any, prefix string) []string {
	var pvs []string
	for key, val := range obj {
		switch t := val.(type) {
		case map[string]any:
			pvs := append(pvs, getPathValues(t, key)...)
			continue
		case []interface{}:
			// Can't handle arrays
			continue
		}

		if prefix != "" {
			key = prefix + "." + key
		}

		pvs = append(pvs, fmt.Sprintf("%s=%v", key, val))
	}

	return pvs
}

func jsonResponse(w http.ResponseWriter, body map[string]any, err error) {
	data := map[string]any{
		"body":   body,
		"status": "ok",
	}

	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		data["status"] = "error"
		data["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		panic(err)
	}
}

func lexString(input []rune, index int) (string, int, error) {
	if index >= len(input) {
		return "", index, nil
	}

	if input[index] == '"' {
		index++
		foundEnd := false

		var s []rune
		for index < len(input) {
			if input[index] == '"' {
				foundEnd = true
				break
			}

			s = append(s, input[index])
			index++
		}

		if !foundEnd {
			return "", index, fmt.Errorf("expected end of quoted string")
		}

		return string(s), index + 1, nil
	}

	// If unquoted, read as much contiguous digits/letters as there are
	var s []rune
	var c rune
	for index < len(input) {
		c = input[index]
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '.') {
			break
		}
		s = append(s, c)
		index++
	}

	if len(s) == 0 {
		return "", index, fmt.Errorf("no string found")
	}

	return string(s), index + 1, nil
}

type queryComparison struct {
	key   []string
	value string
	op    string
}

type query struct {
	ands []queryComparison
}

// E.g. q=a.b:12
func parseQuery(q string) (*query, error) {
	if q == "" {
		return &query{}, nil
	}

	i := 0
	var parsed query
	var qRune = []rune(q)
	for i < len(qRune) {
		// eat whitespace
		for unicode.IsSpace(qRune[i]) {
			i++
		}

		key, nextIndex, err := lexString(qRune, i)
		if err != nil {
			return nil, fmt.Errorf("Expected valid key, got [%s]: `%s`", err, q[nextIndex:])
		}

		// Expect some operator
		if q[nextIndex] != ':' {
			return nil, fmt.Errorf("Expected colon at %d, got: `%s`", nextIndex, q[nextIndex:])
		}
		i = nextIndex + 1

		op := "="
		if q[i] == '>' || q[i] == '<' {
			op = string(q[i])
			i++
		}

		value, nextIndex, err := lexString(qRune, i)
		if err != nil {
			return nil, fmt.Errorf("Expected valid value, got [%s]: `%s`", err, q[nextIndex:])
		}
		i = nextIndex

		argument := queryComparison{key: strings.Split(key, "."), value: value, op: op}
		parsed.ands = append(parsed.ands, argument)
	}

	return &parsed, nil
}

func getPath(doc map[string]any, parts []string) (any, bool) {
	var docSegment any = doc
	for _, part := range parts {
		m, ok := docSegment.(map[string]any)
		if !ok {
			return nil, false
		}

		if docSegment, ok = m[part]; !ok {
			return nil, false
		}
	}

	return docSegment, true
}

func (q query) match(doc map[string]any) bool {
	for _, argument := range q.ands {
		value, ok := getPath(doc, argument.key)
		if !ok {
			return false
		}

		// Handle equality
		if argument.op == "=" {
			match := fmt.Sprintf("%v", value) == argument.value
			if !match {
				return false
			}

			continue
		}

		// Handle <, >
		right, err := strconv.ParseFloat(argument.value, 64)
		if err != nil {
			return false
		}

		var left float64
		switch t := value.(type) {
		case float64:
			left = t
		case float32:
			left = float64(t)
		case uint:
			left = float64(t)
		case uint8:
			left = float64(t)
		case uint16:
			left = float64(t)
		case uint32:
			left = float64(t)
		case uint64:
			left = float64(t)
		case int:
			left = float64(t)
		case int8:
			left = float64(t)
		case int16:
			left = float64(t)
		case int32:
			left = float64(t)
		case int64:
			left = float64(t)
		case string:
			left, err = strconv.ParseFloat(t, 64)
			if err != nil {
				return false
			}
		default:
			return false
		}

		if argument.op == ">" {
			if left <= right {
				return false
			}

			continue
		}

		if left >= right {
			return false
		}
	}

	return true
}

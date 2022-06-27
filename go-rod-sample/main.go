package main

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	page := browser.MustPage("https://pkg.go.dev/time/")

	// wait for footer element is visible (ie, page is loaded)
	page.MustElement(`body > footer`).MustWaitVisible()

	// find and click "Expand All" link
	page.MustElement(`#pkg-examples`).MustClick()

	// retrieve the value of the textarea
	example := page.MustElement(`#example-After textarea`).MustText()

	log.Printf("Go's time.After example:\n%s", example)
}

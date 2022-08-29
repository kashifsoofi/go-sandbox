package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestItems(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Items Suite")
}

type Item struct {
	Id       string
	Quantity int
	SubItems []Item
}

var _ = Describe("Test", func() {
	It("should contain item", func() {
		items := []Item{
			{
				Id:       "1",
				Quantity: 1,
				SubItems: []Item{
					{
						Id:       "1.1",
						Quantity: 11,
					},
				},
			},
			{
				Id:       "2",
				Quantity: 2,
				SubItems: nil,
			},
		}

		Expect(items).To(ContainElement(Item{
			Id:       "1",
			Quantity: 1,
			SubItems: []Item{
				{
					Id:       "1.1",
					Quantity: 11,
				},
			},
		}))
	})
})

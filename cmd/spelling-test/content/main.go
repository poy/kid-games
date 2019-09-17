package main

import (
	"strings"

	"honnef.co/go/js/dom"
)

//go:generate ./generate.sh

func main() {
	d := dom.GetWindow().Document()

	enterButton := d.GetElementByID("enterButton").(*dom.HTMLButtonElement)
	enterButton.Disabled = true

	answerInput := d.GetElementByID("answerInput").(*dom.HTMLInputElement)

	wordInput := d.GetElementByID("wordInput").(*dom.HTMLInputElement)
	wordInput.Style().SetProperty("color", "red", "")
	wordInput.AddEventListener("keyup", true, func(event dom.Event) {
		if answerInput.Value == strings.ToLower(wordInput.Value) {
			wordInput.Style().SetProperty("color", "green", "")
			enterButton.Disabled = false
			return
		}
		wordInput.Style().SetProperty("color", "red", "")
		enterButton.Disabled = true
	})
}

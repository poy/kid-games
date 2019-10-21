package main

import (
	"fmt"
	"strconv"
	"strings"

	"honnef.co/go/js/dom"
)

//go:generate ./generate.sh

func main() {
	d := dom.GetWindow().Document()

	enterButton := d.GetElementByID("enterButton").(*dom.HTMLButtonElement)
	enterButton.Disabled = true

	numberAStr := d.GetElementByID("numberA").(*dom.HTMLInputElement)
	numberBStr := d.GetElementByID("numberB").(*dom.HTMLInputElement)

	answer, err := getAnswer(numberAStr.Value, numberBStr.Value)
	if err != nil {
		fmt.Printf("failed to get answer: %v", err)
		return
	}

	wordInput := d.GetElementByID("wordInput").(*dom.HTMLInputElement)
	wordInput.Style().SetProperty("color", "red", "")
	wordInput.AddEventListener("keyup", true, func(event dom.Event) {
		if answer == strings.ToLower(wordInput.Value) {
			wordInput.Style().SetProperty("color", "green", "")
			enterButton.Disabled = false
			return
		}
		wordInput.Style().SetProperty("color", "red", "")
		enterButton.Disabled = true
	})
}

func getAnswer(strA, strB string) (string, error) {
	a, err := strconv.ParseInt(strA, 10, 64)
	if err != nil {
		return "", err
	}
	b, err := strconv.ParseInt(strB, 10, 64)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(a+b, 10), nil
}

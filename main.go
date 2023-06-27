package main

import (
	"syscall/js"
)

func main() {
	js.Global().Set("reverseText", js.FuncOf(reverseText))
}

func reverseText(this js.Value, args []js.Value) interface{} {
	doc := js.Global().Get("document")
	textbox := doc.Call("getElementById", "textbox")
	label := doc.Call("getElementById", "label")

	text := textbox.Get("value").String()
	reversed := reverseString(text)

	label.Set("textContent", reversed)
	return nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

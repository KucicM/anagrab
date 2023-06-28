package main

import (
	"fmt"
	"syscall/js"
)

const ENTER_KEY_CODE = 13

func main() {
	js.Global().Set("handleKeyPress", js.FuncOf(handleKeyPress))
    
    c := make(chan struct{})
    c <- struct{}{}
}

func handleKeyPress(this js.Value, args []js.Value) interface{} {
    event := args[0]
    keyCode := event.Get("keyCode").Int()
    if keyCode == ENTER_KEY_CODE {
        fmt.Println("Enter")
    }
    fmt.Printf("%d\n", keyCode)
    return nil
}


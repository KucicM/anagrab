package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"syscall/js"
)

const ENTER_KEY_CODE = 13

var wordList []string
var mu = &sync.Mutex{}


func main() {
   js.Global().Set("handleKeyPress", js.FuncOf(handleKeyPress))

   go prefetchWordList()

    <- make(chan struct{})
}

func handleKeyPress(this js.Value, args []js.Value) interface{} {
    event := args[0]
    keyCode := event.Get("keyCode").Int()
    if keyCode == ENTER_KEY_CODE {

        searchWord := getSerachWord(this, args)
        anagrams := findAnagrams(searchWord)
        displayCloud(anagrams)
    }

    return nil
}

func displayCloud(wordList chan string) {
    doc := js.Global().Get("document")
    cloudContainer := doc.Call("getElementById", "result")

    for cloudContainer.Get("firstChild").Truthy() {
		cloudContainer.Call("removeChild", cloudContainer.Get("firstChild"))
	}

    for word := range wordList {
        link := doc.Call("createElement", "a")
        link.Set("textContent", word)

        fontSize := 2.0 + float32(rand() % 5) * 0.5
        opacity := fontSize / 4.5 * 100
		link.Set("style", fmt.Sprintf("font-size: %frem; opacity: %f%%", fontSize, opacity))

        cloudContainer.Call("appendChild", link)
    }

}

func findAnagrams(word string) chan string {
    mu.Lock()
    defer mu.Unlock()

    searchWord := strings.Split(strings.ToLower(word), "")
    sort.Strings(searchWord)

    out := make(chan string)
    go func() {
        defer close(out)
        for i := 0; i < len(wordList); i++ {
            if len(word) != len(wordList[i]) {
                continue
            }

            w := strings.Split(wordList[i], "")
            sort.Strings(w)

            found := true
            for i := 0; i < len(word); i++ {
                if w[i] != searchWord[i] {
                    found = false
                    break
                }
            }

            if found {
                out <- wordList[i]
            }
        }
    }()

    return out
}

func prefetchWordList() {
    mu.Lock()
    defer mu.Unlock()

    words := <- fetchWordList(js.ValueOf(nil), nil)
    wordList = strings.Split(words, "\n")
}

func fetchWordList(this js.Value, args []js.Value) chan string {
    fmt.Println("fetch wordlist")
    ch := make(chan string)

    options := js.Global().Get("Object").New()
	options.Set("headers", js.Global().Get("Object").New())
	options.Get("headers").Set("Accept-Encoding", "gzip")

	respPromise := js.Global().Get("fetch").Invoke("data/word_list.txt", options)
    respPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return args[0].Call("text").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            ch <- args[0].String()
            defer close(ch)
			return nil
		}))
	}))

	return ch
}

func getSerachWord(this js.Value, args []js.Value) string {
    doc := js.Global().Get("document")
    input := doc.Call("getElementById", "search-input")
    return input.Get("value").String()
}

func rand() int {
	return int(js.Global().Get("Math").Call("random").Float() * 100)
}

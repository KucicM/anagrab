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
var mu = &sync.RWMutex{}


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

        for _, anagram := range anagrams {
            fmt.Println(anagram)
        }
    }

    return nil
}

func findAnagrams(word string) []string {
    mu.RLock()
    defer mu.RUnlock()

    searchWord := strings.Split(word, "")
    sort.Strings(searchWord)

    out := make([]string, 0)
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

        if !found {
            continue
        }

        out = append(out, wordList[i])
    }

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

package main

import (
	"fmt"
	"sort"
	"strings"
	"syscall/js"
)

const ENTER_KEY_CODE = 13


func main() {
   js.Global().Set("handleKeyPress", js.FuncOf(handleKeyPress))

    <- make(chan struct{})
}

func handleKeyPress(this js.Value, args []js.Value) interface{} {
    event := args[0]
    keyCode := event.Get("keyCode").Int()
    if keyCode == ENTER_KEY_CODE {

        wordList := getWordList(this, args)

        doc := js.Global().Get("document")
        input := doc.Call("getElementById", "search-input")
        searchWord := input.Get("value").String()
        
        anagrams := findAnagrams(searchWord, wordList)
        for _, anagram := range anagrams {
            fmt.Println(anagram)
        }
    }
    fmt.Printf("%d\n", keyCode)

    return 0
}

func findAnagrams(word string, wordList []string) []string {
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

var wordList []string

func getWordList(this js.Value, args []js.Value) []string {
    if wordList == nil {
        words := <- fetchWordList(this, args)
        wordList = strings.Split(words, "\n")
    }
    return wordList
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

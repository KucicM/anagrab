package main

import (
	"fmt"
	"sort"
	"strings"
	"syscall/js"
	"time"
)

const ENTER_KEY_CODE = 13

var wordList []string
var sortList [][]byte


func main() {
   js.Global().Set("handleKeyPress", js.FuncOf(handleKeyPress))

   prefetchWordList()

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
    lable := doc.Call("getElementById", "not-found-results")

    for cloudContainer.Get("firstChild").Truthy() {
		cloudContainer.Call("removeChild", cloudContainer.Get("firstChild"))
	}

    found := false
    for word := range wordList {
        found = true
        link := doc.Call("createElement", "a")
        link.Set("textContent", word)

        fontSize := 2.0 + float32(rand() % 5) * 0.5
        opacity := fontSize / 4.5 * 100
		link.Set("style", fmt.Sprintf("font-size: %frem; opacity: %f%%", fontSize, opacity))
        link.Set("href", fmt.Sprintf("https://en.wiktionary.org/wiki/%s", word))
        link.Set("target", "_blank")
        link.Set("rel", "noopener noreferrer")

        cloudContainer.Call("appendChild", link)
    }

    if found {
        lable.Get("style").Set("visibility", "hidden")
    } else {
        lable.Get("style").Set("visibility", "visible")
    }

}

func findAnagrams(word string) chan string {
    searchWord := []byte(strings.ToLower(word))
    sort.Slice(searchWord, func(i, j int) bool {return searchWord[i] < searchWord[j]})

    out := make(chan string)
    go func() {
        defer close(out)
        searched := 0
        totalTime := time.Second * 0
        sortTime := time.Second * 0
        defer func() {
            fmt.Printf("search checkd %d words, total %+v sort %+v\n", searched, totalTime, sortTime)
        }()


        for i := 0; i < len(wordList); i++ {
            if len(word) != len(wordList[i]) {
                continue
            }

            searched++
            start := time.Now()

            if sortList[i][0] > searchWord[0] {
                return
            }

            found := true
            for j := 0; j < len(word); j++ {
                if sortList[i][j] != searchWord[j] {
                    found = false
                    break
                }
            }

            totalTime += time.Since(start)

            if found {
                out <- wordList[i]
            }
        }

    }()

    return out
}

func prefetchWordList() {
    words := <- fetchWordList(js.ValueOf(nil), nil)
    wordList = strings.Split(words, "\n")

    sortList = make([][]byte, len(wordList))
    for i := 0; i < len(wordList); i++ {
        lower := strings.ToLower(wordList[i])
        sortList[i] = []byte(lower)
        sort.Slice(sortList[i], func(j int, k int) bool { return sortList[i][j] < sortList[i][k] })
    }

    doc := js.Global().Get("document")
    search := doc.Call("getElementById", "search-input")
    search.Set("placeholder", "")
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

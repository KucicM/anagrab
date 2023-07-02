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
        if len(searchWord) == 0 {
            return nil
        }
        anagrams := findAnagrams(searchWord)
        displayCloud(anagrams, searchWord)
    }

    return nil
}

func displayCloud(wordList chan string, searchWord string) {
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

        editDistance := calculateEditDistance(word, searchWord)
        diff := float32(len(searchWord) - editDistance) / float32(len(searchWord))
        fontSize := 2.0 + diff * 2.5
        opacity := (diff * 0.6 + 0.4) * 100
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
        defer func() {
            fmt.Printf("search checked %d words in total %+v\n", searched, totalTime)
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

func calculateEditDistance(w1, w2 string) int {
    dp := make([]int, len(w2)+1)
    for i := 0; i <= len(w2); i++ {
        prev := dp[0]
        dp[0] = i
        for j := 1; j <= len(w2); j++ {
            tmp := dp[j]
            if i > 0 && j > 0 {
                if w1[i-1] == w2[j-1] {
                    dp[j] = prev
                } else {
                    dp[j] = 1 + min(dp[j], dp[j-1], prev)
                }
            } else {
                dp[j] = j + i
            }
            prev = tmp
        }
    }
    return dp[len(w2)]
}

func min(i, j, k int) int {
    if i <= j && i <= k {
		return i
	} else if j <= i && j <= k {
		return j
	}
	return k
}

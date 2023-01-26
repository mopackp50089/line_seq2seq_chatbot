package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"unicode"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-ego/gse"
)

var (
	seg gse.Segmenter
)

func main() {
	// 加载默认词典
	seg.AlphaNum = true
	err := seg.LoadDict()
	if err != nil {
		log.Println("fail to load default dict", err)
	}
	err = seg.LoadDict("customized_dict.txt")
	if err != nil {
		log.Println("fail to load customized_dict", err)
	}

	// send_hour map
	send_hour := map[int]int{}
	for hour := 0; hour < 24; hour++ {
		send_hour[hour] = 0
	}

	// senders map
	senders := map[string]int{}
	source := "../source/line_jill.txt"
	send_hour, senders, wordsSlice := readline(source, send_hour, senders)

	// content map
	word_counts := map[string]int{}
	for _, word := range wordsSlice {
		rune_word := []rune(word) // for chinese
		if len(rune_word) <= 1 {
			continue
		} else if isDigit(word) {
			continue
		} else if word == "" {
			continue
		} else {
			word_counts[word] = word_counts[word] + 1
		}
	}
	sortKeys := sortByValue(word_counts)

	// TopNumber
	topNumber := 10
	topWords := []string{}
	topCounts := []int{}
	i := -1
	for len(topWords) <= topNumber {
		i += 1
		sortKey := sortKeys[i]
		count := word_counts[sortKey]
		if sortKey == "call" || sortKey == "photo" || sortKey == "vedio" || sortKey == "voice" || sortKey == "sticker" {
			continue
		}
		topWords = append(topWords, sortKey)
		topCounts = append(topCounts, count)
	}

	fmt.Println("topWords", topWords)
	fmt.Println("topCounts", topCounts)

	// summary graph
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{
		Title:    "常用詞Top10",
		Subtitle: "It's extremely easy to use, right?"})

	// Put data into instance
	bar.AddXAxis(topWords).
		AddYAxis("Category A", topCounts)
	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)

}

func sortByValue(word_counts map[string]int) (sortKeys []string) {
	keys := make([]string, 0, len(word_counts))

	for key := range word_counts {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return word_counts[keys[i]] > word_counts[keys[j]]
	})

	return keys
}

func isDigit(str string) bool {
	for _, x := range []rune(str) {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

func readline(path string, send_hour map[int]int, senders map[string]int) (map[int]int, map[string]int, []string) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)
	wordsSlice := []string{}
	for scanner.Scan() {
		// do something with a line
		line := strings.Split(scanner.Text(), "\t")
		if len(line) > 2 {
			timeString, sender, content := line[0], line[1], line[2]
			// timeString
			const timeLayout = "03:04 PM"
			timeStr, err := time.Parse(timeLayout, timeString)
			if err != nil {
				fmt.Println("ERR", err)
			}
			send_hour[timeStr.Hour()] = send_hour[timeStr.Hour()] + 1

			// senders
			if sender != "" {
				senders[sender] = senders[sender] + 1
			}
			// content
			// words := seg.Slice(content)
			words := seg.Slice(content)
			wordsSlice = append(wordsSlice, words...)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return send_hour, senders, wordsSlice

}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
	"flag"
	"sort"
	"strconv"
)

type wordInfo struct{
	timesAppeared int
	seenInLines []int

}

func cleanWord(word string) string{

		//fmt.Printf("Input: %s\n", word)
    runes := []rune(word)
		i := 0
		for{
			if len(runes) >= 1 && !unicode.IsLetter(runes[i]){
				runes = runes[i+1:]
			}else{
				break
			}
		}
		//fmt.Printf("Output: %s\n", string(runes))
    return  string(runes)
}

func saveTargetWords(splittedWords []string, wordsDictionary map[string]wordInfo, lineCounter int) {

	for i := 0; i < len(splittedWords); i++ {
		if (strings.Contains(splittedWords[i], ".") || strings.Contains(splittedWords[i], ":")) && i != len(splittedWords)-1{
			targetWord := splittedWords[i+1]
			if runes := []rune(targetWord); unicode.IsLetter(runes[len(runes)-1]) == false {
				targetWord = targetWord[:len(targetWord)-1]
			}
      	targetWord = cleanWord(targetWord)
				if len(targetWord) >= 3{
					insertInMap(wordsDictionary, targetWord, lineCounter)
				}
		}
	}
}

func insertInMap(wordsDictionary map[string]wordInfo, targetWord string, lineCounter int) {

		targetWord = strings.ToLower(targetWord)
		linesSlice := wordsDictionary[targetWord].seenInLines
		if idx:=sort.SearchInts(linesSlice, lineCounter); idx<len(linesSlice) && len(linesSlice)!=0{
			wordsDictionary[targetWord] = wordInfo{wordsDictionary[targetWord].timesAppeared + 1, linesSlice}

		}else{
			wordsDictionary[targetWord] = wordInfo{wordsDictionary[targetWord].timesAppeared + 1, append(wordsDictionary[targetWord].seenInLines, lineCounter)}
		}
}

func scanAndProcess(scanner *bufio.Scanner, wordsDictionary map[string]wordInfo) {

	var lineCounter int = 1
	for scanner.Scan() {
		line := scanner.Text()
		splittedWords := strings.Fields(line)
		saveTargetWords(splittedWords, wordsDictionary, lineCounter)
		lineCounter++
	}

}

func parseArguments() string{

	filenamePtr := flag.String("file", "", "filename to read")
	flag.Parse()
	if *filenamePtr == ""{
			fmt.Println("Error: at least argument -file is necessary.")
			os.Exit(1)
	}

	return *filenamePtr
}

func prettyPrint(wordsDictionary map[string]wordInfo, filename string){

	for i, v := range wordsDictionary {

        fmt.Printf("%s:%s\n", i, strconv.Itoa(v.timesAppeared))
				for i:=0; i<len(v.seenInLines); i++{
					fmt.Println("\t" + filename + ":" + strconv.Itoa(v.seenInLines[i]))
				}
  }
}

func main() {

	filename := parseArguments()
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: No such file or directory.")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var wordsDictionary map[string]wordInfo
	wordsDictionary = make(map[string]wordInfo)
	scanAndProcess(scanner, wordsDictionary)
	prettyPrint(wordsDictionary, filename)

}

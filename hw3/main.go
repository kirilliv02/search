package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func getInvertedIndex() map[string][]int {
	termDocumentsDict := make(map[string][]int)

	for i := 1; i <= 101; i++ {
		file, _ := os.Open(fmt.Sprintf("lemmas/lemma_%d.txt", i))

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lemma := strings.Split(line, ":")[0]
			termDocumentsDict[lemma] = append(termDocumentsDict[lemma], i)
		}
	}
	return termDocumentsDict
}

func booleanSearch(query string, invertedIndex map[string][]int) []int {

	fields := strings.Split(query, " ")
	res := make([]int, 0)
	operator := ""
	for _, field := range fields {
		if field == "AND" || field == "NOT" || field == "OR" {
			operator = field
			continue
		}

		if len(res) == 0 {
			res = invertedIndex[field]
		} else {
			switch operator {
			case "AND":
				res = and(res, invertedIndex[field])
			case "OR":
				res = or(res, invertedIndex[field])
			case "NOT":
				res = not(res, invertedIndex[field])
			}
		}

	}
	return res
}

func and(a, b []int) []int {
	m := make(map[int]bool)
	var result []int
	for _, val := range a {
		m[val] = true
	}

	for _, val := range b {
		if _, ok := m[val]; ok {
			result = append(result, val)
		}
	}

	return result
}

func or(a, b []int) []int {
	m := make(map[int]bool)
	var result []int

	for _, val := range a {
		m[val] = true
		result = append(result, val)
	}

	for _, val := range b {
		if _, ok := m[val]; !ok {
			result = append(result, val)
		}
	}

	return result
}

func not(a, b []int) []int {
	m := make(map[int]bool)
	var result []int

	for _, val := range b {
		m[val] = true
	}

	for _, val := range a {
		if _, ok := m[val]; !ok {
			result = append(result, val)
		}
	}

	return result
}

func createIndexes() {
	tdDict := getInvertedIndex()

	file, err := os.Create("inverted_index.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for k, v := range tdDict {
		var values []string
		for _, val := range v {
			values = append(values, fmt.Sprintf("%d", val))
		}
		output := fmt.Sprintf("%s %s\n", k, strings.Join(values, " "))
		file.WriteString(output)
	}

	type Index struct {
		Count         int    `json:"count"`
		InvertedArray []int  `json:"inverted_array"`
		Word          string `json:"word"`
	}

	var countInvertedWord []Index
	for k, v := range tdDict {
		countInvertedWord = append(countInvertedWord, Index{
			Count:         len(v),
			InvertedArray: v,
			Word:          k,
		})
	}

	file2, err := os.Create("inverted_index_2.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file2.Close()

	for _, ciw := range countInvertedWord {
		marshal, err := json.Marshal(ciw)
		if err != nil {
			return
		}
		file2.WriteString(string(marshal) + "\n")
	}
}

func main() {
	//createIndexes()

	tdDict := getInvertedIndex()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите запрос: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	searchResults := booleanSearch(query, tdDict)
	fmt.Println("Результаты поиска:")
	fmt.Println(searchResults)
}

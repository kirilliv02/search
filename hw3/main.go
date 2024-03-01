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

func booleanSearch(query []string, invertedIndex map[string][]int) []int {

	res := make([][]int, 0)
	for _, field := range query {
		if field == "AND" || field == "NOT" || field == "OR" {
			first := res[len(res)-2]
			second := res[len(res)-1]
			switch field {
			case "AND":
				res[len(res)-2] = and(first, second)
			case "OR":
				res[len(res)-2] = or(first, second)
			case "NOT":
				res[len(res)-2] = not(first, second)
			}
			res = res[:len(res)-1]
		} else {
			res = append(res, invertedIndex[field])
		}
	}
	return res[0]
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

func infixToPostfix(infix []string) []string {
	var scobs = 0
	var result []string
	var signs []string
	for _, r := range infix {
		if r == "(" {
			scobs++
		} else if r == ")" {
			result = append(result, signs[len(signs)-1])
			signs = signs[:len(signs)-1]
			scobs--
		} else if r == "AND" || r == "OR" || r == "NOT" {
			if len(signs) > 0 {
				if scobs == 0 {
					result = append(result, signs[len(signs)-1])
					signs[len(signs)-1] = r
				} else {
					signs = append(signs, r)
				}

			} else {
				signs = append(signs, r)
			}
		} else {
			result = append(result, r)
		}
	}

	for len(signs) > 0 {
		result = append(result, signs[len(signs)-1])
		signs = signs[:len(signs)-1]
	}

	return result
}

func main() {
	//createIndexes()
	tdDict := getInvertedIndex()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите запрос: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)
	fileds := strings.Fields(query)

	res := make([]string, 0, len(fileds))

	for _, filed := range fileds {
		if filed[0] == '(' {
			for strings.Contains(filed, "(") {
				res = append(res, "(")
				filed = filed[1:]
			}
			res = append(res, filed)

		} else if filed[len(filed)-1] == ')' {
			count := 0
			for strings.Contains(filed, ")") {
				count++
				filed = filed[:len(filed)-1]
			}
			res = append(res, filed)
			for i := 0; i < count; i++ {
				res = append(res, ")")

			}
		} else {
			res = append(res, filed)
		}
	}

	searchResults := booleanSearch(infixToPostfix(res), tdDict)
	fmt.Println("Результаты поиска:")
	fmt.Println(searchResults)
}

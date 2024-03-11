package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// var invertedIdx = getInvertedIndex()
var tfIdfDictsLemmas, idfLemmas = getTFTerms()

func vectorNorm(vec []float64) float64 {
	sum := 0.0
	for _, el := range vec {
		sum += math.Pow(el, 2)
	}
	return math.Pow(sum, 0.5)
}

func getIndex() map[int]string {
	index := make(map[int]string)

	file, err := os.Open("index.txt")
	if err != nil {
		fmt.Println(err)
		return index
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) >= 2 {
			key, err := strconv.Atoi(parts[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
			value := parts[1]
			index[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return index
}

func calculate(term string, documentTokensList []string, documentsCount, documentsWithTermCount int) (float64, float64, float64) {
	tf := float64(strings.Count(strings.Join(documentTokensList, " "), term)) / float64(len(documentTokensList))
	var idf float64
	if documentsWithTermCount == 0 {
		idf = 0
	} else {
		idf = math.Log(float64(documentsCount) / float64(documentsWithTermCount))
	}
	return math.Round(tf*1000000) / 1000000, math.Round(idf*1000000) / 1000000, math.Round(tf*idf*1000000) / 1000000
}

func cosineSimilarity(vec1, vec2 []float64) float64 {
	dot := 0.0
	for i := 0; i < len(vec1); i++ {
		dot += vec1[i] * vec2[i]
	}
	if dot == 0 {
		return 0
	}
	return dot / (vectorNorm(vec1) * vectorNorm(vec2))
}

func search(query string) {
	fmt.Printf("SEARCHING: %s\n", query)
	tokensMap := getLemmas(getTokens(query))

	// Конвертируем значения из карты в срез строк
	var tokens []string
	for _, v := range tokensMap {
		tokens = append(tokens, v...)
	}

	indexDict := getIndex()
	if len(tokens) == 0 {
		fmt.Println("Empty query")
		return
	}
	fmt.Printf("LEMMATIZED: %s\n", strings.Join(tokens, " "))

	queryVector := make([]float64, 0)
	for _, token := range tokens {
		docWithTermsCount := 0
		for _, tfIdfDict := range tfIdfDictsLemmas {
			if _, ok := tfIdfDict[token]; ok {
				docWithTermsCount++
			}
		}
		_, _, tfIdf := calculate(token, tokens, COUNT_DOCUMENTS, docWithTermsCount)
		queryVector = append(queryVector, tfIdf)
	}

	distances := make(map[int]float64)
	for index := 0; index < COUNT_DOCUMENTS; index++ {
		documentVector := make([]float64, 0)
		for _, token := range tokens {
			if tfIdfDict, ok := tfIdfDictsLemmas[index][token]; ok {
				documentVector = append(documentVector, tfIdfDict)
			} else {
				documentVector = append(documentVector, 0.0)
			}
		}
		distances[index] = cosineSimilarity(queryVector, documentVector)
	}

	searchedIndices := make([]int, 0)
	for index := range distances {
		searchedIndices = append(searchedIndices, index)
	}
	sort.Slice(searchedIndices, func(i, j int) bool {
		return distances[searchedIndices[i]] > distances[searchedIndices[j]]
	})

	for _, index := range searchedIndices {
		tfIdf := distances[index]
		if tfIdf < 0.05 {
			continue
		}
		fmt.Printf("Index: %d \n Link: %s \n Cosine: %f \n", index, indexDict[index], tfIdf)
	}
}

func main() {
	var query string
	fmt.Scanln(&query)
	search(query)
}

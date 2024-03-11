package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

var (
	globalIndexToPositions map[string][]int
	indexPositionToURL     map[int]string
)

func main() {
	//var query = "установить NOT освещённость" // освещённость OR установить
	var query string
	fmt.Scanln(&query)
	modifiedSearch(query)
}

func modifiedSearch(queryExpr string) {
	splittedQuery := strings.Split(queryExpr, " ")

	var q query.Query
	if len(splittedQuery) > 1 {
		var prevQuery query.Query
		prevQuery = bleve.NewMatchPhraseQuery(splittedQuery[0])

		for i := 1; i < len(splittedQuery); i++ {

			if i%2 != 0 {
				switch splittedQuery[i] {
				case "AND":
					if i+1 == len(splittedQuery) {
						break
					}
					nextKeyWord := splittedQuery[i+1]
					subQuery := bleve.NewMatchPhraseQuery(nextKeyWord)
					prevQuery = bleve.NewConjunctionQuery(prevQuery, subQuery)

					i++
					continue
				case "OR":
					if i+1 == len(splittedQuery) {
						break
					}
					nextKeyWord := splittedQuery[i+1]
					subQuery := bleve.NewMatchPhraseQuery(nextKeyWord)
					prevQuery = bleve.NewDisjunctionQuery(prevQuery, subQuery)

					i++
					continue
				case "NOT":
					if i+1 == len(splittedQuery) {
						break
					}
					nextKeyWord := splittedQuery[i+1]
					subQuery := bleve.NewMatchPhraseQuery(nextKeyWord)
					tempQuery := bleve.NewBooleanQuery()
					tempQuery.AddMust(prevQuery)
					tempQuery.AddMustNot(subQuery)
					prevQuery = tempQuery

					i++
					continue
				default:
					panic(fmt.Errorf("должен быть логический оператор между поисковыми словами"))
				}
			}
		}

		q = prevQuery
	} else {
		q = bleve.NewMatchPhraseQuery(queryExpr)
	}

	cleanAfteward()

	tree := readIndex()
	search := bleve.NewSearchRequest(q)
	// Параметр отвечающий за то сколько сущностей вернется в ответе
	//search.Size = 10

	res, err := tree.Search(search)
	if err != nil {
		panic(err)
	}

	foundedIndexes := make(map[int]int)
	for _, val := range res.Hits {
		for _, num := range globalIndexToPositions[val.ID] {
			foundedIndexes[num]++
		}
	}

	idx := 0
	for indexPosition, _ := range foundedIndexes {
		idx++
		url := indexPositionToURL[indexPosition]
		fmt.Println(idx, "Подходящий сайт: ", url)
	}

	cleanAfteward()

	//fmt.Println(parsed)
	//fmt.Println(tree)
}

type Index struct {
	Count         int    `json:"count"`
	InvertedArray []int  `json:"inverted_array"`
	Word          string `json:"word"`
}

func cleanAfteward() {
	os.RemoveAll("./index-helper")
}

func readIndex() bleve.Index {
	//fieldMapping :=bleve.NewKeywordFieldMapping()
	mapping := bleve.NewIndexMapping()

	index, err := bleve.New("./index-helper", mapping)
	if err != nil {
		panic(err)
	}
	batch := index.NewBatch()

	file, err := os.Open("inverted_index_2.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	globalIndexToPositions = make(map[string][]int, 10_000)
	for scanner.Scan() {
		idx++
		lineBytes := scanner.Bytes()
		body := Index{}
		err = json.Unmarshal(lineBytes, &body)
		if err != nil {
			panic(err)
		}

		if idx%100 == 0 {
			err = index.Batch(batch)
			if err != nil {
				panic(err)
			}
			batch.Reset()
		}

		globalIndexToPositions[strconv.Itoa(idx)] = body.InvertedArray
		err = batch.Index(strconv.Itoa(idx), body)
		if err != nil {
			panic(err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при сканировании файла:", err)
	}

	indexPositionToURL = make(map[int]string, 101)
	file2, err := os.Open("index.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		panic(err)
	}
	defer file2.Close()

	scanner2 := bufio.NewScanner(file2)
	for scanner2.Scan() {
		line := scanner2.Text()
		splitted := strings.Split(line, " ")

		num, err := strconv.Atoi(splitted[0])
		if err != nil {
			panic(err)
		}
		indexPositionToURL[num] = splitted[1]
	}

	return index
}

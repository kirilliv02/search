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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var (
	globalIndexToPositions map[string][]int
	indexPositionToURL     map[int]string
	tree                   bleve.Index
)

func main() {
	cleanAfteward()
	defer cleanAfteward()
	tree = readIndex()

	engine := html.New("./demo/templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Обработчик для отдачи страницы с формой
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	// Обработчик для обработки данных формы
	app.Post("/submit", func(c *fiber.Ctx) error {
		data := c.FormValue("data")
		results := modifiedSearch(data)
		return c.Render("results", fiber.Map{
			"Results": results,
		})
	})

	app.Listen(":8080")

}

type Result struct {
	Url string
}

func modifiedSearch(queryExpr string) []Result {
	queryExpr = strings.ToLower(queryExpr)
	splittedQuery := strings.Split(queryExpr, " ")

	var q query.Query
	if len(splittedQuery) > 1 {
		var prevQuery query.Query
		prevQuery = bleve.NewMatchPhraseQuery(splittedQuery[0])

		for i := 1; i < len(splittedQuery); i++ {

			if i%2 != 0 {
				switch splittedQuery[i] {
				case "and":
					if i+1 == len(splittedQuery) {
						break
					}
					nextKeyWord := splittedQuery[i+1]
					subQuery := bleve.NewMatchPhraseQuery(nextKeyWord)
					prevQuery = bleve.NewConjunctionQuery(prevQuery, subQuery)

					i++
					continue
				case "or":
					if i+1 == len(splittedQuery) {
						break
					}
					nextKeyWord := splittedQuery[i+1]
					subQuery := bleve.NewMatchPhraseQuery(nextKeyWord)
					prevQuery = bleve.NewDisjunctionQuery(prevQuery, subQuery)

					i++
					continue
				case "not":
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
	resp := make([]Result, 0, len(foundedIndexes))
	for indexPosition, _ := range foundedIndexes {
		idx++
		resp = append(resp, Result{Url: indexPositionToURL[indexPosition]})
	}

	//fmt.Println(parsed)
	//fmt.Println(tree)

	return resp
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

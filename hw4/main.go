package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	DIRECTORY       = "pages"
	COUNT_DOCUMENTS = 101
)

func countTF(terms []string) map[string]float64 {
	tfDict := make(map[string]float64)
	for _, term := range terms {
		tfDict[term]++
	}
	for k, v := range tfDict {
		tfDict[k] = v / float64(len(terms))
	}
	return tfDict
}

func countIDF(terms []string, termsInDocuments [][]string) map[string]float64 {
	idfDict := make(map[string]float64)
	for _, term := range terms {
		countDocWithTerm := 0
		for _, termsInDocument := range termsInDocuments {
			for _, t := range termsInDocument {
				if t == term {
					countDocWithTerm++
					break
				}
			}
		}
		idfDict[term] = math.Log(float64(COUNT_DOCUMENTS) / float64(countDocWithTerm))
	}
	return idfDict
}

func countTFIDF(tfDict map[string]float64, idfDict map[string]float64) map[string]float64 {
	tfIDFDict := make(map[string]float64)
	for term, tfValue := range tfDict {
		tfIDFDict[term] = tfValue * idfDict[term]
	}
	return tfIDFDict
}

func getTFTerms(kind string) ([]map[string]float64, map[string]float64) {
	var termsOverall []string
	var termsInDocuments [][]string
	tfDocuments := make([]map[string]float64, COUNT_DOCUMENTS)

	idx := -1

	err := filepath.Walk(DIRECTORY, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".txt") && strings.HasPrefix(strings.ToLower(info.Name()), kind) {
			idx++
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			terms := strings.Split(string(data), "\n")
			termsInDocuments = append(termsInDocuments, terms)
			tfDocuments[idx] = countTF(terms)
			termsOverall = append(termsOverall, terms...)
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	termsOverall = removeDuplicates(termsOverall)
	idfTerms := countIDF(termsOverall, termsInDocuments)

	var tfIDFDicts []map[string]float64
	for _, tfDocument := range tfDocuments {
		tfIDFTerms := countTFIDF(tfDocument, idfTerms)
		tfIDFDicts = append(tfIDFDicts, tfIDFTerms)
	}

	return tfIDFDicts, idfTerms
}

func writingFiles(tfIDFDicts []map[string]float64, idfTerms map[string]float64, kind string) {
	pathResultBegin := fmt.Sprintf("tf_idf_%s/tf_idf_%s_", kind, kind)
	for i, tfIDFDict := range tfIDFDicts {
		numFile := fmt.Sprintf("%03d", i)
		pathResult := fmt.Sprintf("%s%s.txt", pathResultBegin, numFile)

		file, err := os.Create(pathResult)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer file.Close()

		for k, v := range tfIDFDict {
			line := fmt.Sprintf("%s %f %f\n", k, idfTerms[k], v)
			_, err := file.WriteString(line)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func main() {
	tfIDFDictsTokens, idfTokens := getTFTerms("tokens")
	writingFiles(tfIDFDictsTokens, idfTokens, "tokens")

	tfIDFDictsLemmas, idfLemmas := getTFTerms("lemmas")
	writingFiles(tfIDFDictsLemmas, idfLemmas, "lemmas")
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kljensen/snowball/russian"
	"golang.org/x/text/unicode/norm"
)

var stopwords = map[string]struct{}{
	"и":       {},
	"в":       {},
	"во":      {},
	"не":      {},
	"что":     {},
	"он":      {},
	"на":      {},
	"я":       {},
	"с":       {},
	"со":      {},
	"как":     {},
	"а":       {},
	"то":      {},
	"все":     {},
	"она":     {},
	"так":     {},
	"его":     {},
	"но":      {},
	"да":      {},
	"ты":      {},
	"к":       {},
	"у":       {},
	"же":      {},
	"вы":      {},
	"за":      {},
	"бы":      {},
	"по":      {},
	"только":  {},
	"ее":      {},
	"мне":     {},
	"было":    {},
	"вот":     {},
	"от":      {},
	"меня":    {},
	"еще":     {},
	"нет":     {},
	"о":       {},
	"из":      {},
	"ему":     {},
	"теперь":  {},
	"когда":   {},
	"даже":    {},
	"ну":      {},
	"вдруг":   {},
	"ли":      {},
	"если":    {},
	"уже":     {},
	"или":     {},
	"ни":      {},
	"быть":    {},
	"был":     {},
	"него":    {},
	"до":      {},
	"вас":     {},
	"нибудь":  {},
	"опять":   {},
	"уж":      {},
	"вам":     {},
	"ведь":    {},
	"там":     {},
	"потом":   {},
	"себя":    {},
	"ничего":  {},
	"ей":      {},
	"может":   {},
	"они":     {},
	"тут":     {},
	"где":     {},
	"есть":    {},
	"надо":    {},
	"ней":     {},
	"для":     {},
	"мы":      {},
	"тебя":    {},
	"их":      {},
	"чем":     {},
	"была":    {},
	"сам":     {},
	"чтоб":    {},
	"без":     {},
	"будто":   {},
	"чего":    {},
	"раз":     {},
	"тоже":    {},
	"себе":    {},
	"под":     {},
	"будет":   {},
	"ж":       {},
	"тогда":   {},
	"кто":     {},
	"этот":    {},
	"того":    {},
	"потому":  {},
	"этого":   {},
	"какой":   {},
	"совсем":  {},
	"ним":     {},
	"здесь":   {},
	"этом":    {},
	"один":    {},
	"почти":   {},
	"мой":     {},
	"тем":     {},
	"чтобы":   {},
	"нее":     {},
	"сейчас":  {},
	"были":    {},
	"куда":    {},
	"зачем":   {},
	"всех":    {},
	"никогда": {},
	"можно":   {},
	"при":     {},
	"наконец": {},
	"два":     {},
	"об":      {},
	"другой":  {},
	"хоть":    {},
	"после":   {},
	"над":     {},
	"больше":  {},
	"тот":     {},
	"через":   {},
	"эти":     {},
	"нас":     {},
	"про":     {},
	"всего":   {},
	"них":     {},
	"какая":   {},
	"много":   {},
	"разве":   {},
	"три":     {},
	"эту":     {},
	"моя":     {},
	"впрочем": {},
	"хорошо":  {},
	"свою":    {},
	"этой":    {},
	"перед":   {},
	"иногда":  {},
	"лучше":   {},
	"чуть":    {},
	"том":     {},
	"нельзя":  {},
	"такой":   {},
	"им":      {},
	"более":   {},
	"всегда":  {},
	"конечно": {},
	"всю":     {},
	"между":   {},
}

const directory = "Downloads"
const directory2 = "Tokens"
const directory3 = "Lemmas"

func getTokens(s string) []string {
	re := regexp.MustCompile("[А-Яа-яёЁ]+")
	cleanWords := re.FindAllString(s, -1)
	var tokens []string
	for _, w := range cleanWords {
		w = strings.ToLower(w)
		if _, ok := stopwords[w]; !ok && !containsNumber(w) {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

func containsNumber(s string) bool {
	re := regexp.MustCompile("[0-9]")
	return re.MatchString(s)
}

func getLemmas(tokens []string) map[string][]string {
	lemmas := make(map[string][]string)
	for _, token := range tokens {
		lemma := russian.Stem(token, false)
		if _, ok := lemmas[lemma]; !ok {
			lemmas[lemma] = []string{}
		}
		lemmas[lemma] = append(lemmas[lemma], token)
	}
	return lemmas
}

func processFile(filePath string) ([]string, map[string][]string) {
	htmlText, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlText)))
	if err != nil {
		log.Fatal(err)
	}
	var texts []string
	doc.Find("body").Children().Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(norm.NFC.String(s.Text()))
		texts = append(texts, text)
	})
	result := strings.Join(texts, " ")

	tokens := getTokens(result)
	uniqueTokens := uniqueSlice(tokens)
	lemmas := getLemmas(uniqueTokens)

	// Сохранение токенов в файл
	tokensFilePath := filepath.Join(directory2, "tokens_"+filepath.Base(filePath))
	tokensFile, err := os.Create(tokensFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer tokensFile.Close()
	for _, token := range tokens {
		_, err := tokensFile.WriteString(token + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Сохранение лемм в файл
	lemmasFilePath := filepath.Join(directory3, "lemmas_"+filepath.Base(filePath))
	lemmasFile, err := os.Create(lemmasFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer lemmasFile.Close()
	for lemma, words := range lemmas {
		_, err := lemmasFile.WriteString(lemma + ": " + strings.Join(words, " ") + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	return uniqueTokens, lemmas
}

func uniqueSlice(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func processAllFiles() {
	for _, file := range getFiles() {
		tokens, lemmas := processFile(file)
		fmt.Println("Tokens for file", file, ":", tokens)
		fmt.Println("Lemmas for file", file, ":")
		for lemma, words := range lemmas {
			fmt.Println(lemma, ":", words)
		}
	}
}

func getFiles() []string {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".txt") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

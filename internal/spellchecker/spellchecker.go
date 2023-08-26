package spellchecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	responseTime  = 1
	spellCheckURL = "https://speller.yandex.net/services/spellservice.json/checkText"
)

type SpellResult struct {
	Code    int      `json:"code"`
	Pos     int      `json:"pos"`
	Row     int      `json:"row"`
	Col     int      `json:"col"`
	Len     int      `json:"len"`
	Word    string   `json:"word"`
	Suggest []string `json:"s"`
}

func Handle(text string) (string, error) {
	client := http.Client{Timeout: time.Second * responseTime}

	requestURL := fmt.Sprintf("%s?text=%s", spellCheckURL, strings.ReplaceAll(text, " ", "+"))
	resp, err := client.Post(requestURL, "text/plain", nil)
	if err != nil {
		return text, fmt.Errorf("network error: %v", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return text, err
	}

	var spellResult []SpellResult
	err = json.Unmarshal(buf.Bytes(), &spellResult)
	if spellResult != nil && err != nil {
		return text, err
	}

	return restoreText(text, &spellResult), nil
}

func restoreText(input string, spellResult *[]SpellResult) string {
	tmpInput := []rune(input)
	var res []rune

	var i int
	for _, spellWord := range *spellResult {
		if i >= len(input) {
			break
		}
		substr := tmpInput[i:spellWord.Pos]
		res = append(res, substr...)
		res = append(res, []rune(spellWord.Suggest[0])...)

		i = spellWord.Pos + spellWord.Len
	}

	if i < len(tmpInput) {
		res = append(res, tmpInput[i:]...)
	}
	return string(res)
}

package main

import (
	"os"
	"log"
	"bytes"
	"io"
	"io/ioutil"
	"bufio"
	"encoding/json"
	"net/http"
	"github.com/joho/godotenv"
	"golang.org/x/text/transform"
)

var APIKEY string

func init(){
	err := godotenv.Load()
    if err != nil {
		panic(err)
	}
	APIKEY =  os.Getenv("APIKEY")
}

type CodicEngine struct {
	Text string `json:"text"`
}

type CodicWordsCandidates struct {
	Text string `json:"text"`
}

type CodicWords struct {
	Successful bool `json:"successful"`
	Text string `json:"text"`
	Translated_text string `json:"translated_text"`
	Candidates []CodicWordsCandidates `json:"candidates"`
}

type CodicRes struct {
	Successful bool `json:"successful"`
	Text string `json:"text"`
	Translated_text string `json:"translated_text"`
	Words []CodicWords `json:"words"`
}

func transformEncoding( rawReader io.Reader, trans transform.Transformer) (string, error) {
    ret, err := ioutil.ReadAll(transform.NewReader(rawReader, trans))
    if err == nil {
        return string(ret), nil
    } else {
        return "", err
    }
}

func codic(text string) (*[]CodicRes, error){
	client := &http.Client{}
	url := "https://api.codic.jp/v1/engine/translate.json"
	method := "POST"
	jsonBytes, err := json.Marshal(CodicEngine{Text: text})
	if err != nil{
		return nil, err
	}
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil{
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + APIKEY)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Origi", "*")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		return nil, err
	}
	defer res.Body.Close()

	data := new([]CodicRes)
	if err = json.Unmarshal([]byte(body), data); err != nil {
		return nil, err
	}

	return data, nil
}

func main(){
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()
	data, err := codic(text)
	if err != nil{
		log.Println(err)
	} else {
		log.Println(data)
	}
}
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/text/encoding/charmap"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName xml.Name `xml:"Valute"`
	Name    string   `xml:"Name"`
}

func main() {
	response, err := http.Get("http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=11/11/2020")
	if err != nil {
		log.Fatalf("failed send request: %v", err)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("failed read response body: %v", err)
	}

	fmt.Println(string(b))

	var valCurs ValCurs
	d := xml.NewDecoder(bytes.NewReader(b))
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	err = d.Decode(&valCurs)
	if err != nil {
		log.Fatalf("failed Unmarshal body: %v", err)
	}

	for _, valute := range valCurs.Valutes {
		fmt.Println(valute)
	}
}

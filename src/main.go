package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/encoding/charmap"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName xml.Name `xml:"Valute"`
	Name    string   `xml:"Name"`
	Value   string   `xml:"Value"`
}

const (
	DDMMYYYY = "02/01/2006"
)

func main() {
	currentDate := time.Now().Local().AddDate(0, 0, -90)
	fmt.Println(currentDate.Format(DDMMYYYY))

	response, err := http.Get(fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=%v", currentDate.Format(DDMMYYYY)))
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

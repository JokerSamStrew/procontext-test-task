package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"log"

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

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unknown charset: %s", charset)
	}
}

func cbrRequest(date time.Time) ([]byte, error) {
	const (
		Attempts = 10
		Delay    = 5
	)

	var err error
	for i := 0; i < Attempts; i++ {
		url := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=%v", date.Format(DDMMYYYY))
		response, err := http.Get(url)
		log.Printf("Send request: %v", url)

		if err != nil {
			return nil, fmt.Errorf("failed send request: %v", err)
		}

		b, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed read response body: %v", err)
		}

		if response.StatusCode == http.StatusOK {
			return b, nil
		}

		log.Printf("Attempt failed. Bad response Body: \n%v\n Try again after %v seconds...\n", string(b), Delay)
		time.Sleep(Delay * time.Second)
	}

	return nil, err
}

func RetrieveValCurs(date time.Time) ([]Valute, error) {
	b, err := cbrRequest(date)
	if err != nil {
		return nil, err
	}

	var valCurs ValCurs
	d := xml.NewDecoder(bytes.NewReader(b))
	d.CharsetReader = charsetReader

	err = d.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("failed parse body: %v", err)
	}

	return valCurs.Valutes, nil
}

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func cursValToFloat(value string) float64 {
	curs, err := strconv.ParseFloat(strings.Replace(value, ",", ".", -1), 32)
	if err != nil {
		log.Fatalf("failed parse curs val: %v", err)
	}

	return curs
}

func main() {
	const (
		DateRange = 90
	)

	for i := 0; i < DateRange; i++ {
		currentDate := time.Now().Local().AddDate(0, 0, -i)
		valutes, err := RetrieveValCurs(currentDate)
		if err != nil {
			log.Fatalf("failed retrieve val curs: %v", err)
		}

		for _, valute := range valutes {
			fmt.Println(cursValToFloat(valute.Value))
		}
	}
}

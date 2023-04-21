package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	currentDate := time.Now().Local().AddDate(0, 0, -90)

	valutes, err := RetrieveValCurs(currentDate)
	if err != nil {
		log.Fatalf("failed retrieve val curs: %v", err)
	}

	for _, valute := range valutes {
		fmt.Println(valute)
	}
}

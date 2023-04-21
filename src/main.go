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

type CursData struct {
	date time.Time
	curs float64
}

const (
	DateRange = 90
)

func retrieveCursData() map[string][]*CursData {
	data := make(map[string][]*CursData)
	for i := 0; i < DateRange; i++ {
		currentDate := time.Now().Local().AddDate(0, 0, -i)
		valutes, err := RetrieveValCurs(currentDate)
		if err != nil {
			log.Fatalf("failed retrieve val curs: %v", err)
		}

		for _, valute := range valutes {
			data[valute.Name] = append(data[valute.Name], &CursData{date: currentDate, curs: cursValToFloat(valute.Value)})
		}
	}

	return data
}

func main() {
	type ResultStruct struct {
		maxCurs *CursData
		minCurs *CursData
		sumCurs float64
		count   float64
	}

	result := make(map[string]*ResultStruct)

	for name, cursData := range retrieveCursData() {
		fmt.Printf("%v", name)

		for _, val := range cursData {
			fmt.Printf(" {%v %v}", val.date.Format("02-01-2006"), val.curs)

			if result[name] == nil {
				result[name] = &ResultStruct{maxCurs: val, minCurs: val, sumCurs: val.curs, count: 1}
				continue
			}

			if result[name].maxCurs.curs < val.curs {
				result[name].maxCurs = val
			}

			if result[name].minCurs.curs > val.curs {
				result[name].minCurs = val
			}

			result[name].sumCurs += val.curs
			result[name].count += 1
		}

		fmt.Println()
	}

	for name, val := range result {
		fmt.Printf("\n%v: \n\tmax curs: %v date: %v\n\tmin curs: %v date: %v\n\tavgCurs: %v\n",
			name,
			val.maxCurs.curs, val.maxCurs.date.Format("02-01-2006"),
			val.minCurs.curs, val.minCurs.date.Format("02-01-2006"),
			val.sumCurs/val.count,
		)
	}
}

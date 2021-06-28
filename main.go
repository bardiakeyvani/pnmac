package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {

	isWeatherTable := flag.Bool("weather", false, "if table contains weather data")
	isSoccerTable := flag.Bool("soccer", false, "if table contains weather data")
	tablePath := flag.String("table", "", "path to where table file is located")

	flag.Parse()

	if len(*tablePath) == 0 {
		fmt.Print("please provide comple path to table file. example: ./tables/w_data.dat")
		return
	}
	if !*isWeatherTable && !*isSoccerTable {
		fmt.Print("please provide table type. example --weather, --soccer")
		return
	}
	if *isWeatherTable {
		processWeatherTable(*tablePath)
		return
	}

	processSoccerTable(*tablePath)
}

func processWeatherTable(tablePath string) {
	rows, err := readWeatherTable(tablePath)
	if err != nil {
		log.Fatalf("Failed to read weather table. error: %s", err.Error())
	}

	dayWithSmallestSpread, err := getRowWithSmallestTemperatureSpread(rows)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("smallest temperature spread is on day: ", dayWithSmallestSpread[0])
}

func processSoccerTable(filepath string) {
	rows, err := readSoccerTable(filepath)
	if err != nil {
		log.Fatalf("Failed to read soccer table. error: %s", err.Error())
	}

	teamWithSmallestScoreDifference, err := findTeamWithSmallestScoreDifference(rows)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Team with smallest goal difference is: ", teamWithSmallestScoreDifference[1])
}

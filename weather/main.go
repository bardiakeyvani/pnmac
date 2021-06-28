package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var skipWeatherRowsRegex = regexp.MustCompile("<pre>|</pre>|MMU|Dy|mo")
var tableRowRegex = regexp.MustCompile(`\s+(\d+)\s+(\d+)\*?\s+(\d+)\*?\s+(\d+)\s+(\d+)?\s+(\d+\.\d)\s+(\d+\.\d{2})\s+([A-Z]+)?\s+(\d{3})\s+(\d+\.\d{1})\s(\d{3})\s+(\d+)\*?\s+(\d+.\d)\s+(\d+)\s+(\d+)\s+(\d+.\d)$`)

func main() {

	tablePath := flag.String("table", "", "path to where table file is located")

	flag.Parse()

	if len(*tablePath) == 0 {
		fmt.Print("please provide comple path to table file. example: ./tables/w_data.dat")
		return
	}

	processWeatherTable(*tablePath)

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

func getRowWithSmallestTemperatureSpread(rows [][]string) ([]string, error) {
	if len(rows) == 0 || len(rows[0]) != 16 {
		return nil, errors.New("not a valid table")
	}
	smallestSpread := 0
	smallestSpreadRow := 0
	for i, row := range rows {
		maxTemp, err := strconv.ParseInt(row[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MxT: %s on row %d", row[1], i)
		}
		minTemp, _ := strconv.ParseInt(row[2], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MnT: %s on row %d", row[2], i)
		}
		spread := int(maxTemp) - int(minTemp)

		log.Printf("day: %s, maxT: %d, minT: %d, spread: %d", row[0], maxTemp, minTemp, spread)
		if i == 0 {
			smallestSpread = spread
			continue
		}

		if smallestSpread > spread {
			smallestSpread = spread
			smallestSpreadRow = i
		}
	}

	return rows[smallestSpreadRow], nil
}
func readWeatherTable(fileName string) ([][]string, error) {
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %s", err.Error())
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	table := [][]string{}
	for scanner.Scan() {
		row := scanner.Text()
		// skip row if it matches
		if skipWeatherRowsRegex.MatchString(row) {
			continue
		}

		if len(row) == 0 {
			continue
		}

		if len(row) < 84 {
			return nil, fmt.Errorf("row length is less than 84, actual: %d", len(row))
		}

		table = append(table, getWeatherRowColumns(row))
	}

	return table, nil
}

func getWeatherRowColumns(row string) []string {
	// remove all *
	columns := tableRowRegex.FindAllStringSubmatch(row, -1)
	if len(columns) != 1 && len(columns[0]) != 17 {
		log.Fatalf("not a valid row. count: %d", len(columns[0]))
	}
	return []string{
		columns[0][1],
		columns[0][2],
		columns[0][3],
		columns[0][4],
		columns[0][5],
		columns[0][6],
		columns[0][7],
		columns[0][8],
		columns[0][9],
		columns[0][10],
		columns[0][11],
		columns[0][12],
		columns[0][13],
		columns[0][14],
		columns[0][15],
		columns[0][16],
	}

}

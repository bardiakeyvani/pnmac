package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var skipWeatherRowsRegex = regexp.MustCompile("<pre>|</pre>|MMU|Dy|mo")

func getRowWithSmallestTemperatureSpread(rows [][]string) ([]string, error) {
	if len(rows) == 0 || len(rows[0]) != 18 {
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
	row = strings.ReplaceAll(row, "*", " ")
	return []string{
		strings.Replace(row[0:4], " ", "", -1),
		strings.Replace(row[5:10], " ", "", -1),
		strings.Replace(row[12:17], " ", "", -1),
		strings.Replace(row[17:22], " ", "", -1),
		strings.Replace(row[23:29], " ", "", -1),
		strings.Replace(row[30:35], " ", "", -1),
		strings.Replace(row[35:40], " ", "", -1),
		strings.Replace(row[41:46], " ", "", -1),
		strings.Replace(row[47:53], " ", "", -1),
		strings.Replace(row[47:53], " ", "", -1),
		strings.Replace(row[54:58], " ", "", -1),
		strings.Replace(row[59:63], " ", "", -1),
		strings.Replace(row[64:67], " ", "", -1),
		strings.Replace(row[68:71], " ", "", -1),
		strings.Replace(row[72:76], " ", "", -1),
		strings.Replace(row[76:80], " ", "", -1),
		strings.Replace(row[80:83], " ", "", -1),
		strings.Replace(row[83:], " ", "", -1),
	}

}

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

var skipSoccerRowsRegex = regexp.MustCompile("<pre>|</pre>|Team|--")
var rowRegex = regexp.MustCompile(`\s+(\d+)\.\s+([A-Za-z_]+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+-\s+(\d+)\s+(\d+)$`)

func main() {

	tablePath := flag.String("table", "", "path to where table file is located")

	flag.Parse()

	if len(*tablePath) == 0 {
		fmt.Print("please provide comple path to table file. example: ./tables/w_data.dat")
		return
	}
	processSoccerTable(*tablePath)
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

func readSoccerTable(fileName string) ([][]string, error) {
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %s", err.Error())
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	table := [][]string{}
	// loop as long as there is more line
	for scanner.Scan() {
		row := scanner.Text()
		if skipSoccerRowsRegex.MatchString(row) {
			continue
		}

		if len(row) == 0 {
			continue
		}
		if len(row) < 56 {
			return nil, fmt.Errorf("not a valid table. Row length is less than 58, actual: %d", len(row))
		}

		table = append(table, getSoccerRowColumns(row))
	}

	return table, nil
}

func getSoccerRowColumns(row string) []string {
	// replace all with "-" with " "
	columns := rowRegex.FindAllStringSubmatch(row, -1)
	if len(columns) != 1 && len(columns[0]) != 10 {
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
	}
}

func findTeamWithSmallestScoreDifference(rows [][]string) ([]string, error) {
	if len(rows) == 0 || len(rows[0]) != 9 {
		return nil, errors.New("not a valid table soccer")
	}
	smallestScoreDiff := 0
	smallestScoreDiffRowIndex := 0
	for i, row := range rows {
		forGoals, err := strconv.ParseInt(row[6], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse F: %s on row %d", row[6], i)
		}
		againstGoals, _ := strconv.ParseInt(row[7], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse A: %s on row %d", row[7], i)
		}
		diffGoal := Abs(int(forGoals) - int(againstGoals))

		log.Printf("Team: %s, For Goals: %d, Against Goals: %d, Diff Goals: %d", row[1], forGoals, againstGoals, diffGoal)
		if i == 0 {
			smallestScoreDiff = diffGoal
			continue
		}

		if smallestScoreDiff > diffGoal {
			smallestScoreDiff = diffGoal
			smallestScoreDiffRowIndex = i
		}
	}

	return rows[smallestScoreDiffRowIndex], nil
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

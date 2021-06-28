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

var skipSoccerRowsRegex = regexp.MustCompile("<pre>|</pre>|Team|--")

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
	row = strings.ReplaceAll(row, "-", " ")
	return []string{
		strings.Replace(row[0:5], " ", "", -1),
		strings.Replace(row[7:23], " ", "", -1),
		strings.Replace(row[23:28], " ", "", -1),
		strings.Replace(row[29:33], " ", "", -1),
		strings.Replace(row[33:37], " ", "", -1),
		strings.Replace(row[37:42], " ", "", -1),
		strings.Replace(row[37:42], " ", "", -1),
		strings.Replace(row[42:47], " ", "", -1),
		strings.Replace(row[42:47], " ", "", -1),
		strings.Replace(row[48:55], " ", "", -1),
		strings.Replace(row[55:], " ", "", -1),
	}
}

func findTeamWithSmallestScoreDifference(rows [][]string) ([]string, error) {
	if len(rows) == 0 || len(rows[0]) != 11 {
		return nil, errors.New("not a valid table soccer")
	}
	smallestScoreDiff := 0
	smallestScoreDiffRowIndex := 0
	for i, row := range rows {
		forGoals, err := strconv.ParseInt(row[8], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse F: %s on row %d", row[8], i)
		}
		againstGoals, _ := strconv.ParseInt(row[9], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse A: %s on row %d", row[9], i)
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

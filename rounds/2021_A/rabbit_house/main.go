package main

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

type testCase struct {
	rows int
	cols int
	grid [][]int
}

type testCaseOrErr struct {
	testCase
	err error
}

func main() {

}

// -------- Input reading -------- //

func NewTestCase(rows, cols int, heights [][]int) testCase {
	return testCase{
		rows,
		cols,
		heights,
	}
}

func NewTestCaseOrErr(rows, cols int, heights [][]int, err error) testCaseOrErr {
	return testCaseOrErr{
		NewTestCase(rows, cols, heights),
		err,
	}
}

func parseIntFields(line string) (ints []int, err error) {
	for _, field := range strings.Fields(line) {
		convField, err := strconv.Atoi(field)
		if err != nil {
			return []int{}, err
		}
		ints = append(ints, convField)
	}
	return
}

// Go
func LoadCsvDataToChannel(in io.Reader) <-chan testCaseOrErr {
	out := make(chan PointOrErr)
	go func() {
		defer close(out)
		reader := csv.NewReader(in)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				return
			}
			if err != nil {
				out <- PointOrErr{Err: err}
				return
			}
			point, err := recordToPoint(record)
			if err != nil {
				out <- PointOrErr{Err: err}
				return
			}
			out <- PointOrErr{Point: point}
		}
	}()
	return out
}

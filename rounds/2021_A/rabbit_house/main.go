package main

import (
	"bufio"
	"fmt"
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

func newTestCase(rows, cols int, heights [][]int) testCase {
	return testCase{
		rows,
		cols,
		heights,
	}
}

func newTestCaseOrErr(rows, cols int, grid [][]int, err error) testCaseOrErr {
	return testCaseOrErr{
		newTestCase(rows, cols, grid),
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

func parseRowAndColNum(reader bufio.Reader) (row, col int, err error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	
	intFields, err := parseIntFields(line)
	if err != nil {
		return
	}

	if len(intFields) != 2 {
		err = fmt.Errorf("number of int fields in first line of test case not equal to 2")
		return
	}

	row = intFields[0]
	col = intFields[1]
	
	return
}

func parseGrid(rows int, cols int, reader bufio.Reader) ([][]int, error) {
	var grid [rows][cols]int
	
	for i := 0; i < rows; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return grid, err
		}
	
		row, err := parseIntFields(line)
		if err != nil {
			return grid, err
		}
	
		grid[i] = row
	}

	return grid, nil
}

func loadTestCasesToChannel(reader bufio.Reader) <-chan testCaseOrErr {
	out := make(chan testCaseOrErr)
	
	go func() {
		defer close(out)
		
		numberOfTestCases, err := reader.ReadString('\n')
		if err != nil {
			out <- testCaseOrErr{err: err}
			
			return
		}

		for i := 0; i < numberOfTestCases; i++ {
			rows, cols, err := parseRowAndColNum(reader)
			if err != nil {
				out <- testCaseOrErr{err: err}
				
				return
			}

			grid, err := parseGrid(rows, cols, reader)
			if err != nil {
				out <- testCaseOrErr{err: err}
				
				return
			}

			out <- newTestCaseOrErr(rows, cols, grid)
		}

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				out <- testCaseOrErr{err: err}
				return
			}

			, err := recordToPoint(record)
			if err != nil {
				out <- PointOrErr{Err: err}
				return
			}
			out <- PointOrErr{Point: point}
		}
	}()

	return out
}

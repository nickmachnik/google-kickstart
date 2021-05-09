package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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
	reader := bufio.NewReader(os.Stdin)
	testCases := loadTestCasesToChannel(reader)
	for test := range testCases {
		if test.err != nil {
			log.Fatal(test.err)
		}
		fmt.Println(test)
	}
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

func parseIntsFromNextLine(reader *bufio.Reader) (ints []int, err error) {
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return
	}

	return parseIntFields(line)
}

func parseRowAndColNum(reader *bufio.Reader) (row, col int, err error) {
	intFields, err := parseIntsFromNextLine(reader)
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

func parseNumTestCases(reader *bufio.Reader) (numTestCases int, err error) {
	firstLineInts, err := parseIntsFromNextLine(reader)
	if err != nil {
		return
	}

	if len(firstLineInts) != 1 {
		err = fmt.Errorf("unexpected number of ints in test case number definition")

		return
	}

	numTestCases = firstLineInts[0]

	return
}

func parseGrid(rows int, cols int, reader *bufio.Reader) ([][]int, error) {
	grid := make([][]int, rows)

	for i := 0; i < rows; i++ {
		row, err := parseIntsFromNextLine(reader)
		if err != nil {
			return grid, err
		}

		grid[i] = row
	}

	return grid, nil
}

func loadTestCasesToChannel(reader *bufio.Reader) <-chan testCaseOrErr {
	out := make(chan testCaseOrErr)

	go func() {
		defer close(out)

		numberOfTestCases, err := parseNumTestCases(reader)
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

			out <- newTestCaseOrErr(rows, cols, grid, err)
		}
	}()

	return out
}

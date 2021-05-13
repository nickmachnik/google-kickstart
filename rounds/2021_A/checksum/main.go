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

func main() {
	reader := bufio.NewReader(os.Stdin)
	testCases := loadTestCasesToChannel(reader)
	var testIx int
	for test := range testCases {
		testIx++
		if test.err != nil {
			log.Fatal(test.err)
		}
		// fmt.Printf("Case #%d: %d\n", testIx, numAdditions)
	}
}

// -------- Input reading -------- //

type testCase struct {
	n            int
	boolMatrix   [][]int
	costMatrix   [][]int
	rowChecksums []int
	colChecksums []int
}

type testCaseOrErr struct {
	testCase testCase
	err      error
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
			out <- parseTestCase(reader)
		}
	}()

	return out
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

func parseTestCase(reader *bufio.Reader) testCaseOrErr {
	dim, err := parseDimensions(reader)
	if err != nil {
		return testCaseOrErr{err: err}
	}

	boolMatrix, err := parseMatrix(dim, reader)
	if err != nil {
		return testCaseOrErr{err: err}
	}

	costMatrix, err := parseMatrix(dim, reader)
	if err != nil {
		return testCaseOrErr{err: err}
	}

	rowChecksum, err := parseChecksum(dim, reader)
	if err != nil {
		return testCaseOrErr{err: err}
	}

	colChecksum, err := parseChecksum(dim, reader)
	if err != nil {
		return testCaseOrErr{err: err}
	}

	return testCaseOrErr{
		testCase{
			dim,
			boolMatrix,
			costMatrix,
			rowChecksum,
			colChecksum,
		},
		nil,
	}
}

func parseDimensions(reader *bufio.Reader) (dim int, err error) {
	intFields, err := parseIntsFromNextLine(reader)
	if err != nil {
		return
	}

	if len(intFields) != 1 {
		err = fmt.Errorf("number of int fields in first line of test case not equal to 1")
		return
	}

	dim = intFields[0]

	return
}

func parseMatrix(n int, reader *bufio.Reader) ([][]int, error) {
	matrix := make([][]int, n)

	for i := 0; i < n; i++ {
		row, err := parseIntsFromNextLine(reader)
		if err != nil {
			return matrix, err
		}

		matrix[i] = row
	}

	return matrix, nil
}

func parseChecksum(n int, reader *bufio.Reader) (checksum []int, err error) {
	checksum, err = parseIntsFromNextLine(reader)
	if err != nil {
		return
	}

	if len(checksum) != n {
		err = fmt.Errorf("expected %d vals in checksum, got %d", n, len(checksum))
		return
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

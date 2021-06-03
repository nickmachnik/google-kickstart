package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
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
		fmt.Printf("Case #%d: %d\n", testIx, countMaxProductiveInputs(&test.testCase))
	}
}

// -------- Test solving --------- //

func countMaxProductiveInputs(test *testCase) (count int) {
	maxI := int(math.Sqrt(2.0 * float64(test.g)))
	for i := 1; i <= maxI; i++ {
		if hasPositiveIntegerRoot(i, test.g) {
			count++
		}
	}

	return
}

func hasPositiveIntegerRoot(i, g int) bool {
	root := float64(2*g+i-i*i) / float64(2*i)

	return isIntegral(root) && root >= 0
}

func isIntegral(val float64) bool {
	return val == float64(int(val))
}

// -------- Input reading -------- //

type testCase struct {
	g int
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
	g, err := parseG(reader)
	if err != nil {
		return testCaseOrErr{err: err}
	}

	return testCaseOrErr{
		testCase{g},
		nil,
	}
}

func parseG(reader *bufio.Reader) (g int, err error) {
	intFields, err := parseIntsFromNextLine(reader)
	if err != nil {
		return
	}

	if len(intFields) != 1 {
		err = fmt.Errorf("number of int fields in first line of test case not equal to 1")
		return
	}

	g = intFields[0]

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

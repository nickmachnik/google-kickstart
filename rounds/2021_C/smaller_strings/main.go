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

// alphabetIndex returns the 0 based index of a lowercase letter in the english alphabet
// no bound checks!
func alphabetIndex(c rune) int {
	return int(c) - 97
}

func countSmallerPalindromes(test *testCase) (res int) {
	half := int(math.Ceil(float64(test.n) / 2))
	s := []rune(test.s)

	for i := 0; i < half; i++ {
		res += alphabetIndex(s[i]) * int(math.Pow(float64(test.k), float64(half-i-1)))
	}

	// at least one smaller, none larger
	var anyLarger, anySmaller bool
	for i := half; i < test.n; i++ {
		if s[test.n-i-1] > s[i] {
			anyLarger = true
		} else if s[test.n-i-1] < s[i] {
			anySmaller = true
		}
		fmt.Println(res, s[test.n-i-1], s[i], anyLarger, anySmaller)
	}
	if !anyLarger && anySmaller {
		res++
	}

	return
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	testCases := loadTestCasesToChannel(reader)
	var testIx int
	for test := range testCases {
		testIx++
		if test.err != nil {
			log.Fatal(test.err)
		}
		// fmt.Println(test.testCase)
		fmt.Printf("Case #%d: %d\n", testIx, countSmallerPalindromes(&test.testCase)%int(math.Pow(10, 9)+7))
	}
}

// -------- Input reading -------- //

type testCase struct {
	n int
	k int
	s string
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
	n, k, err := parseNK(reader)
	if err != nil {
		return testCaseOrErr{err: err}
	}

	s, err := parseStringFromNextLine(reader)
	if err != nil {
		return testCaseOrErr{err: err}
	}

	return testCaseOrErr{
		testCase{n, k, s},
		nil,
	}
}

func parseNK(reader *bufio.Reader) (n, k int, err error) {
	intFields, err := parseIntsFromNextLine(reader)
	if err != nil {
		return
	}

	if len(intFields) != 2 {
		err = fmt.Errorf("number of int fields in first line of test case not equal to 2")
		return
	}

	n = intFields[0]
	k = intFields[1]

	return
}

func parseStringFromNextLine(reader *bufio.Reader) (s string, err error) {
	s, err = reader.ReadString('\n')
	if err == io.EOF {
		return s, nil
	}

	s = strings.TrimSuffix(s, "\n")
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

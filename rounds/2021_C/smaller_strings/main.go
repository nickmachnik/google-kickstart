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

func countSmallerPalindromes(s []rune, variablePos, alphabetSize int) (count int, inputIsPalindrome bool) {
	n := len(s)

	if n == 0 {
		return 0, true
	}

	smallerLetters := alphabetIndex(s[0])
	if n == 1 {
		return smallerLetters, true
	}

	innerPalindromes := math.Pow(float64(alphabetSize), float64(variablePos))
	innerCount, innerIsPalindrome := countSmallerPalindromes(s[1:n-1], variablePos-1, alphabetSize)
	inputIsPalindrome = innerIsPalindrome && s[0] == s[n-1]
	if innerIsPalindrome && s[0] < s[n-1] {
		innerCount++
	}
	return smallerLetters*int(innerPalindromes) + innerCount, inputIsPalindrome

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
		s := []rune(test.testCase.s)
		variablePos := int(math.Ceil(float64(test.testCase.n)/float64(2))) - 1
		palCount, _ := countSmallerPalindromes(s, variablePos, test.testCase.k)
		fmt.Printf("Case #%d: %d\n", testIx, palCount%int(math.Pow(10, 9)+7))
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

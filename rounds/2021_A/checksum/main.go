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

type boolMatrix [][]int
type costMatrix [][]int

type entryValue struct {
	known bool
	value bool
}

type entry struct {
	row   int
	col   int
	cost  int
	value entryValue
}

type entries map[*entry]struct{}

type rowOrColumn struct {
	entries        entries
	observedParity bool
	expectedParity bool
}

type dependencyGroup struct {
	resolvableEntries map[entry]struct{}
	rows              map[int]*rowOrColumn
	cols              map[int]*rowOrColumn
}

func newDependencyGroup() *dependencyGroup {
	return &dependencyGroup{
		resolvableEntries: make(map[entry]struct{}),
		rows:              make(map[int]*rowOrColumn),
		cols:              make(map[int]*rowOrColumn),
	}
}

func buildDependencyGroup(startingRowIndex int, test *testCase) (g *dependencyGroup, coveredRows map[int]struct{}) {
	coveredRows = make(map[int]struct{})
	g = newDependencyGroup()
	rowsToAdd := make(map[int]struct{})
	colsToAdd := make(map[int]struct{})
	rowsToAdd[startingRowIndex] = struct{}{}

	for len(rowsToAdd) != 0 || len(colsToAdd) != 0 {
		for rowIndex := range rowsToAdd {
			for _, colIndex := range g.addRow(rowIndex, test) {
				colsToAdd[colIndex] = struct{}{}
			}
			coveredRows[rowIndex] = struct{}{}
			delete(rowsToAdd, rowIndex)
		}

		for colIndex := range colsToAdd {
			for _, rowIndex := range g.addCol(colIndex, test) {
				rowsToAdd[rowIndex] = struct{}{}
			}
			delete(colsToAdd, colIndex)
		}
	}

	return
}

func (g *dependencyGroup) addRow(rowIndex int, test *testCase) (colsToAdd []int) {
	newRow := rowOrColumn{
		expectedParity: test.rowChecksums[rowIndex] == 1,
		entries:        make(entries),
	}

	for colIndex, value := range test.boolMatrix[rowIndex] {
		if value < 0 {
			if _, ok := g.cols[colIndex]; !ok {
				colsToAdd = append(colsToAdd, colIndex)
			}
			newRow.entries[newUnknownEntry(rowIndex, colIndex, test.costMatrix[rowIndex][colIndex])] = struct{}{}
		} else if value == 1 {
			newRow.flipParity()
		}
	}

	g.rows[rowIndex] = &newRow
	if len(newRow.entries) == 1 {
		var e *entry
		for e = range newRow.entries {
			break
		}
		g.resolvableEntries[*e] = struct{}{}
	}

	return
}

func (g *dependencyGroup) addCol(colIndex int, test *testCase) (rowsToAdd []int) {
	newCol := rowOrColumn{
		expectedParity: test.colChecksums[colIndex] == 1,
		entries:        make(entries),
	}

	for rowIndex := 0; rowIndex < test.n; rowIndex++ {
		value := test.boolMatrix[rowIndex][colIndex]
		if value < 0 {
			if _, ok := g.rows[rowIndex]; !ok {
				rowsToAdd = append(rowsToAdd, rowIndex)
			}
			newCol.entries[newUnknownEntry(rowIndex, colIndex, test.costMatrix[rowIndex][colIndex])] = struct{}{}
		} else if value == 1 {
			newCol.flipParity()
		}
	}

	g.cols[colIndex] = &newCol
	if len(newCol.entries) == 1 {
		var e *entry
		for e = range newCol.entries {
			break
		}
		g.resolvableEntries[*e] = struct{}{}
	}

	return
}

func (r *rowOrColumn) flipParity() {
	r.observedParity = !r.observedParity
}

func newUnknownEntry(row int, col int, cost int) *entry {
	return &entry{row, col, cost, entryValue{false, false}}
}

func (g *dependencyGroup) resolve() (cost int) {

	return
}

func (g *dependencyGroup) resolveResolvable() {
	for e := range g.resolvableEntries {
		e.value.known = true

		if len(g.rows[e.row].entries) == 1 {
			e.value.value = g.rows[e.row].expectedParity != g.rows[e.row].observedParity
			g.rows

		} else {
			e.value.value = g.cols[e.col].expectedParity != g.cols[e.col].observedParity
		}

		delete(g.resolvableEntries, e)
	}
}

func resolveMatrix(test *testCase) (cost int) {
	solvedRows := make(map[int]struct{})
	for rowIndex := 0; rowIndex < test.n; rowIndex++ {
		if _, ok := solvedRows[rowIndex]; ok {
			continue
		}
		for colIndex := 0; colIndex < test.n; colIndex++ {
			if test.boolMatrix[rowIndex][colIndex] < 0 {
				group, groupRows := buildDependencyGroup(rowIndex, test)
				for rowIndex := range groupRows {
					solvedRows[rowIndex] = struct{}{}
				}
				fmt.Println(group)
				cost += group.resolve()
				break
			}
		}
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
		fmt.Println(test.testCase)
		resolveMatrix(&test.testCase)
		fmt.Println()
	}
}

// -------- Input reading -------- //

type testCase struct {
	n            int
	boolMatrix   boolMatrix
	costMatrix   costMatrix
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

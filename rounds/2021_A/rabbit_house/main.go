package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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
	var testIx int
	for test := range testCases {
		testIx++
		if test.err != nil {
			log.Fatal(test.err)
		}
		for _, row := range test.grid {
			fmt.Println(row)
		}
		fmt.Println("------")
		numAdditions := makeRabbitHouseSafe(&test.testCase)
		for _, row := range test.grid {
			fmt.Println(row)
		}
		fmt.Printf("Case #%d: %d\n", testIx, numAdditions)
	}
}

func makeRabbitHouseSafe(house *testCase) (totalAddedHeight int) {
	flatGrid := flattenGrid(house.rows, house.cols, house.grid)
	sort.Slice(flatGrid, func(i, j int) bool {
		return flatGrid[i].height > flatGrid[j].height
	})

	auxillaryQueue := make(chan gridLocation, len(flatGrid))
	var lastHeight int
	for _, loc := range flatGrid {
		if loc.height != lastHeight {
			numQueued := len(auxillaryQueue)
			for i := 0; i < numQueued; i++ {
				queuedLoc := <-auxillaryQueue
				totalAddedHeight += secureLocation(&queuedLoc, house, auxillaryQueue)
			}
		}
		totalAddedHeight += secureLocation(&loc, house, auxillaryQueue)
		lastHeight = loc.height
	}

	return
}

func secureLocation(loc *gridLocation, house *testCase, auxillaryQueue chan gridLocation) (addedHeight int) {
	if loc.height != house.grid[loc.row][loc.col] {
		return
	}
	addedHeight += adjustNeighborHeights(loc, house, auxillaryQueue)

	return
}

func adjustNeighborHeights(loc *gridLocation, house *testCase, auxillaryQueue chan gridLocation) (addedHeight int) {
	for _, neighbor := range getNeighborLocs(loc, house) {
		heightDiff := loc.height - neighbor.height
		if heightDiff > 1 {
			addedHeight += heightDiff - 1
			house.grid[neighbor.row][neighbor.col] = loc.height - 1
			neighbor.height = loc.height - 1
			auxillaryQueue <- neighbor
		}
	}

	return
}

func getNeighborLocs(loc *gridLocation, house *testCase) (neighbors []gridLocation) {
	if loc.row > 0 {
		neighbors = append(neighbors, gridLocation{loc.row - 1, loc.col, house.grid[loc.row-1][loc.col]})
	}
	if loc.col < house.cols-1 {
		neighbors = append(neighbors, gridLocation{loc.row, loc.col + 1, house.grid[loc.row][loc.col+1]})
	}
	if loc.row < house.rows-1 {
		neighbors = append(neighbors, gridLocation{loc.row + 1, loc.col, house.grid[loc.row+1][loc.col]})
	}
	if loc.col > 0 {
		neighbors = append(neighbors, gridLocation{loc.row, loc.col - 1, house.grid[loc.row][loc.col-1]})
	}

	return
}

type gridLocation struct {
	row, col, height int
}

func flattenGrid(rows, cols int, grid [][]int) []gridLocation {
	flatGrid := make([]gridLocation, 0, rows*cols)
	for rowIx, row := range grid {
		for colIx, height := range row {
			flatGrid = append(flatGrid, gridLocation{rowIx, colIx, height})
		}
	}
	return flatGrid
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

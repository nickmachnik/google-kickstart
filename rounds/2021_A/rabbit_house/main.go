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
	var testIx int
	for test := range testCases {
		testIx++
		if test.err != nil {
			log.Fatal(test.err)
		}
		numAdditions := makeRabbitHouseSafe(&test.testCase)
		fmt.Printf("Case #%d: %d\n", testIx, numAdditions)
	}
}

func makeRabbitHouseSafe(house *testCase) (totalHeightIncrease int) {
	buckets := newHeightBuckets(&house.grid)
	for {
		if buckets.maxHeight == 0 {
			break
		}
		totalHeightIncrease += secureNextLocation(buckets, house)
	}

	return
}

func secureNextLocation(buckets *heightBuckets, house *testCase) (addedHeight int) {
	loc := buckets.getLocationAtMaxHeight()
	locHeight := getLocationHeight(loc, house)
	defer buckets.removeLocation(locHeight, loc)

	for _, neighbor := range getNeighborLocations(loc, house) {
		neighborHeight := getLocationHeight(neighbor, house)
		heightDiff := locHeight - neighborHeight
		if heightDiff > 1 {
			addedHeight += heightDiff - 1
			buckets.insertLocation(locHeight-1, neighbor)
			setLocationHeight(locHeight-1, neighbor, house)
			buckets.removeLocation(neighborHeight, neighbor)
		}
	}

	return
}

func setLocationHeight(height int, loc location, house *testCase) {
	house.grid[loc.row][loc.col] = height
}

func getLocationHeight(loc location, house *testCase) (height int) {
	return house.grid[loc.row][loc.col]
}

func getNeighborLocations(loc location, house *testCase) (neighbors []location) {
	if loc.row > 0 {
		neighbors = append(neighbors, location{loc.row - 1, loc.col})
	}
	if loc.col < house.cols-1 {
		neighbors = append(neighbors, location{loc.row, loc.col + 1})
	}
	if loc.row < house.rows-1 {
		neighbors = append(neighbors, location{loc.row + 1, loc.col})
	}
	if loc.col > 0 {
		neighbors = append(neighbors, location{loc.row, loc.col - 1})
	}

	return
}

type location struct {
	row, col int
}

type heightBuckets struct {
	buckets   map[int]map[location]struct{}
	maxHeight int
}

func (b *heightBuckets) getLocationAtMaxHeight() location {
	loc, err := b.getLocationAtHeight(b.maxHeight)
	if err != nil {
		log.Fatal(err)
	}

	return loc
}

func (b *heightBuckets) getLocationAtHeight(height int) (loc location, err error) {
	for loc = range b.buckets[height] {

		return loc, err
	}

	return loc, fmt.Errorf("no location found at height: %d", height)
}

func (b *heightBuckets) insertLocation(height int, loc location) {
	if _, ok := b.buckets[height]; !ok {
		b.buckets[height] = map[location]struct{}{}
	}
	b.buckets[height][loc] = struct{}{}
	if height > b.maxHeight {
		b.maxHeight = height
	}
}

func (b *heightBuckets) removeLocation(height int, loc location) {
	delete(b.buckets[height], loc)
	if len(b.buckets[height]) == 0 {
		delete(b.buckets, height)
	}
	if height == b.maxHeight {
		b.decreaseMaxHeight()
	}
}

func (b *heightBuckets) decreaseMaxHeight() {
	if len(b.buckets) == 0 {
		b.maxHeight = 0

		return
	}
	for {
		if _, ok := b.buckets[b.maxHeight]; !ok {
			b.maxHeight--
		} else {
			break
		}
	}
}

func newHeightBuckets(grid *[][]int) *heightBuckets {
	ret := heightBuckets{
		buckets:   make(map[int]map[location]struct{}),
		maxHeight: 0,
	}
	for rowIx, row := range *grid {
		for colIx, height := range row {
			ret.insertLocation(height, location{rowIx, colIx})
		}
	}
	return &ret
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

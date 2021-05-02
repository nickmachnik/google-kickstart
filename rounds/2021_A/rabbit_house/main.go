package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type testCase struct {
	rows int
	cols int
	grid [][]int
	err  error
}

type testCaseOrErr struct {
	testCase
	err error
}

func main() {

}

// -------- Input reading -------- //

// Go

func NewTestCaseOrErr(rawLines []string) (res testCaseOrErr) {
	strings.Split() rawLines[0]
	strconv.ParseInt() 
	return
}

func recordToPoint(record []string) (p Point, err error) {
	if len(record) != 2 {
		err = fmt.Errorf("Records must have two columns")
		return
	}
	if p.X, err = strconv.ParseFloat(record[0], 64); err != nil {
		return
	}
	if p.Y, err = strconv.ParseFloat(record[1], 64); err != nil {
		return
	}
	return
}

// Go
func LoadCsvDataToChannel(in io.Reader) <-chan PointOrErr {
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

package pkg

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// CSVSource is a struct that holds the data from the CSV file
type CSVSource struct {
	Header []string
	Rows   [][]string
}

// NewCSVSource creates a new CSVSource struct
func NewCSVSource() *CSVSource {
	return &CSVSource{}
}

// ReadCSV reads the CSV file and stores the data in the CSVSource struct
func (c *CSVSource) ReadCSV(filename string) error {
	csvFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	// reader.Comma = ';'

	var header []string
	var rows [][]string

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading CSV file: %v\n", err)
			return err
		}

		if header == nil {
			header = row
			continue
		}

		rows = append(rows, row)
	}

	c.Header = header
	c.Rows = rows

	return nil
}

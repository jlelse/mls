package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	// Check if the correct number of command-line arguments is provided
	if len(os.Args) < 6 {
		fmt.Println("Usage: go run main.go <input CSV file> <output CSV file> <radio> <mcc> <network> <min samples> <output columns>")
		return
	}

	// Get the command-line arguments
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	desiredRadio := os.Args[3]
	desiredMCC := os.Args[4]
	desiredNetwork := os.Args[5]
	desiredMinSamples := os.Args[6]
	desiredColumns := os.Args[7:]

	// Open the input CSV file
	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening the input file:", err)
		return
	}
	defer input.Close()

	// Create the output CSV file
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating the output file:", err)
		return
	}
	defer output.Close()

	// CSV reader and writer initialization
	reader := csv.NewReader(input)
	writer := csv.NewWriter(output)

	// Read the header row to get the column indices
	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading CSV header:", err)
		return
	}

	radioIndex := findColumnIndex(header, "radio")
	mccIndex := findColumnIndex(header, "mcc")
	netIndex := findColumnIndex(header, "net")
	samplesIndex := findColumnIndex(header, "samples")

	// Find the indices of the desired columns
	desiredIndices := make([]int, len(desiredColumns))
	for i, column := range desiredColumns {
		desiredIndices[i] = findColumnIndex(header, column)
		if desiredIndices[i] == -1 {
			fmt.Println("Column not found:", column)
			return
		}
	}

	// Write the header row to the output file
	err = writer.Write(desiredColumns)
	if err != nil {
		fmt.Println("Error writing the header row:", err)
		return
	}

	// Read the CSV data line by line
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data:", err)
			return
		}

		// Check the filter criteria
		radio := row[radioIndex]
		mcc := row[mccIndex]
		net := row[netIndex]
		samples := row[samplesIndex]

		if radio == desiredRadio && mcc == desiredMCC && net == desiredNetwork && samples >= desiredMinSamples {
			// Extract the desired columns based on indices
			filteredRow := make([]string, len(desiredIndices))
			for i, index := range desiredIndices {
				filteredRow[i] = row[index]
			}

			// Write the row to the output file
			err := writer.Write(filteredRow)
			if err != nil {
				fmt.Println("Error writing the row:", err)
				return
			}
		}
	}

	// Write CSV writer data to the file
	writer.Flush()

	if err := writer.Error(); err != nil {
		fmt.Println("Error writing the output file:", err)
		return
	}

	fmt.Println("Filtering completed.")
}

// Helper function to find the index of a column name in the header row
func findColumnIndex(header []string, columnName string) int {
	for i, name := range header {
		if name == columnName {
			return i
		}
	}
	return -1
}

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	leafletVersion       = "1.9.4"
	markerClusterVersion = "1.4.1"
	htmlTemplateStart    = `<!DOCTYPE html>
<html>
<head>
	<title>Location Map</title>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/leaflet@%s/dist/leaflet.css" />
	<link rel="stylesheet" href="https://unpkg.com/leaflet.markercluster@%s/dist/MarkerCluster.css" />
	<link rel="stylesheet" href="https://unpkg.com/leaflet.markercluster@%s/dist/MarkerCluster.Default.css" />
	<style>
		#map {
			height: 800px;
		}
	</style>
</head>
<body>
	<div id="map"></div>
	<script src="https://cdn.jsdelivr.net/npm/leaflet@%s/dist/leaflet.js"></script>
	<script src="https://unpkg.com/leaflet.markercluster@%s/dist/leaflet.markercluster.js"></script>
	<script>
		var map = L.map('map').setView([0, 0], 2);

		L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
			attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors',
			maxZoom: 18,
		}).addTo(map);

		var markers = L.markerClusterGroup();

`
	htmlTemplateEnd = `
		map.addLayer(markers);
	</script>
</body>
</html>`
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input CSV file (needs lon,lat,created,updated,cell)> <output HTML file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Open the CSV file
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the CSV records
	reader := csv.NewReader(file)
	records, err := readCSVRecords(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Generate the HTML
	html := generateHTML(records)

	// Save the HTML to a file
	err = writeHTMLToFile(outputFile, html)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Map generated successfully and saved to %s.\n", outputFile)
}

func readCSVRecords(reader *csv.Reader) ([][]string, error) {
	reader.Comma = ',' // Set the delimiter to comma explicitly
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Remove the header row
	if len(records) > 0 {
		records = records[1:]
	}

	return records, nil
}

func generateHTML(records [][]string) string {
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString(fmt.Sprintf(htmlTemplateStart, leafletVersion, markerClusterVersion, markerClusterVersion, leafletVersion, markerClusterVersion))

	for _, record := range records {
		// Parse longitude and latitude
		lon, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatalf("Error parsing longitude: %v\nLongitude value: %s\n", err, record[0])
		}
		lat, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatalf("Error parsing latitude: %v\nLatitude value: %s\n", err, record[1])
		}

		// Get created and updated timestamps
		created := record[2]
		updated := record[3]
		cell := record[4]

		htmlBuilder.WriteString(fmt.Sprintf(`
		var marker = L.marker([%f, %f]);
		marker.bindTooltip("Created: %s<br />Updated: %s<br />Cell: %s");
		markers.addLayer(marker);
		`, lat, lon, formatTimestamp(created), formatTimestamp(updated), cell))
	}

	htmlBuilder.WriteString(htmlTemplateEnd)
	return htmlBuilder.String()
}

func formatTimestamp(timestamp string) string {
	// Parse the Unix timestamp string to an integer
	unixTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "" // Return an empty string if parsing fails
	}

	// Convert the Unix timestamp to a time.Time value
	timeValue := time.Unix(unixTime, 0)

	// Format the time
	isoFormatted := timeValue.Format(time.RFC3339)

	return isoFormatted
}

func writeHTMLToFile(filename string, html string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(html)
	if err != nil {
		return err
	}

	return nil
}

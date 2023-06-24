# Cell Location Map using Mozilla Location Service

This repository provides code to display the cell locations using the Mozilla Location Service. The code extracts cell data from a CSV export file and generates an interactive map with markers representing the cell locations.

## Prerequisites

- Go programming language
- Access to a [cell export file from the Mozilla Location Service](https://location.services.mozilla.com/downloads)

## Installation

1. Clone the repository:

   ```shell
   git clone https://git.jlel.se/jlelse/mls.git
   ```

2. Change to the project directory:

   ```shell
   cd mls
   ```

## Usage

### Step 1: Filter Cell Data

To filter the cell data and generate an output file, run the following command:

```shell
go run ./filter/main.go <input CSV file> <output CSV file> <radio> <mcc> <network> <min samples> <output columns>
```

- `<input CSV file>`: Path to the input CSV file containing the cell data.
- `<output CSV file>`: Path to the output CSV file where the filtered data will be saved.
- `<radio>`: The desired radio type (e.g., LTE).
- `<mcc>`: The desired Mobile Country Code (e.g., 262 for Germany).
- `<network>`: The desired network code (e.g., 3 for o2-de).
- `<min samples>`: The minimum number of samples required for a cell to be included.
- `<output columns>`: Space-separated list of columns to include in the output.

Example:

```shell
go run ./filter/main.go ./MLS-full-cell-export.csv ./output.csv LTE 262 3 100 lon lat created updated cell
```

### Step 2: Generate Map

To generate an HTML map with the filtered cell data, run the following command:

```shell
go run ./web/main.go <input CSV file> <output HTML file>
```

- `<input CSV file>`: Path to the input CSV file containing the filtered cell data with the columns lon,lat,created,updated,cell.
- `<output HTML file>`: Path to the output HTML file where the map will be saved.

Example:

```shell
go run ./web/main.go ./output.csv ./map.html
```

## Acknowledgments

The code in this repository is based on the [Mozilla Location Service](https://location.services.mozilla.com/) and leverages the following open-source libraries:

- [Leaflet](https://leafletjs.com/) (version 1.9.4): An open-source JavaScript library for interactive maps.
- [Leaflet.markercluster](https://github.com/Leaflet/Leaflet.markercluster) (version 1.4.1): A plugin for Leaflet that provides clustering functionality for markers.

The idea for this project was inspired by the [CellInspector Map](https://mls.maxomagier.de/map.html).

## License

This project is licensed under the [MIT License](LICENSE).
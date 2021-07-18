package main

import (
	"bluesky.com/greenhouse-gas-emissions/connections"
	"bluesky.com/greenhouse-gas-emissions/models"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	// uncomment to create tables in db
	// connections.Db.AutoMigrate(&models.Country{}, &models.EmissionCategory{}, &models.CountryEmission{})
}

func main() {

	inventory_data := getInventoryDataFromCsv()
	populateData(inventory_data)

}

func populateData(inventory_data []models.GreenhouseGas) {

	country_data := make([]models.Country, 1)
	emission_category_data := make([]models.EmissionCategory, 1)
	country_emission_data := make([]models.CountryEmission, 1)

	unique_country := make(map[string]bool)
	unique_emission_category := make(map[string]bool)

	for _, row := range inventory_data {
		country := new(models.Country)
		emission_category := new(models.EmissionCategory)

		// country start
		_, is_country_present := unique_country[row.CountryOrArea]
		if is_country_present == false {
			country.Name = row.CountryOrArea
			if country_data[0].Name == "" {
				country_data[0] = *country
			} else {
				country_data = append(country_data, *country)
			}
			unique_country[country.Name] = true
		}
		// country end

		// emission category start
		category := row.Category
		_, is_emission_category_present := unique_emission_category[category]
		if is_emission_category_present == false {
			emission_category.Name = category
			if emission_category_data[0].Name == "" {
				emission_category_data[0] = *emission_category
			} else {
				emission_category_data = append(emission_category_data, *emission_category)
			}
			unique_emission_category[category] = true
		}
		// emission category end

	}

	// country emission start
	country_map := populateCountry(country_data)
	emission_category_map := populateEmissionCategory(emission_category_data)

	for _, row := range inventory_data {
		country_emission := new(models.CountryEmission)
		country_emission.CountryID = country_map[row.CountryOrArea]
		country_emission.EmissionCategoryID = emission_category_map[row.Category]
		country_emission.Year = row.Year
		country_emission.Value = row.Value

		if country_emission_data[0].CountryID == 0 {
			country_emission_data[0] = *country_emission
		} else {
			country_emission_data = append(country_emission_data, *country_emission)
		}

	}
	fmt.Println("Inserting country emissions")
	connections.Db.CreateInBatches(country_emission_data, 9000)
	// country emission end

}

func populateCountry(country []models.Country) map[string]int {
	fmt.Println("Inserting countries")
	connections.Db.CreateInBatches(country, 9000)
	country_map := make(map[string]int)
	for _, x := range country {
		country_map[x.Name] = x.ID
	}
	return country_map
}

func populateEmissionCategory(emission_category []models.EmissionCategory) map[string]int {
	fmt.Println("Inserting emission categories")
	connections.Db.CreateInBatches(emission_category, 9000)
	emission_category_map := make(map[string]int)
	for _, x := range emission_category {
		emission_category_map[x.Name] = x.ID
	}
	return emission_category_map
}

func getInventoryDataFromCsv() []models.GreenhouseGas {
	fmt.Println("Reading csv data")

	file, err := os.Open("greenhouse_gas_inventory_data_data.csv")
	if err != nil {
		log.Fatalf("Error %s opening file greenhouse_gas_inventory_data_data.csv: ", err)
	}
	// close csv file
	defer file.Close()

	reader := bufio.NewReader(file)

	// read the first row before the loop as it contains the header line
	// rows are new line separated
	reader.ReadString('\n')

	inventory_data := make([]models.GreenhouseGas, 1)
	for {
		// read one row from the file
		// rows are new line separated
		row, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		row = string(row[:len(row)-1])

		// columns are comma separated
		columns := strings.Split(row, ",")

		greenhouse_gas := new(models.GreenhouseGas)
		greenhouse_gas.CountryOrArea = columns[0]
		greenhouse_gas.Year, _ = strconv.Atoi(columns[1])
		greenhouse_gas.Value, _ = strconv.ParseFloat(columns[2], 64)
		greenhouse_gas.Category = columns[3]

		if inventory_data[0].CountryOrArea == "" {
			inventory_data[0] = *greenhouse_gas
		} else {
			inventory_data = append(inventory_data, *greenhouse_gas)
		}

	}

	return inventory_data
}

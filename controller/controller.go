package controller

import (
	"bluesky.com/greenhouse-gas-emissions/connections"
	"bluesky.com/greenhouse-gas-emissions/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func Countries(c *gin.Context) {
	var countries_api_response []models.CountriesApiResponse
	connections.Db.Table("countries").Select("countries.id, countries.name, MIN(year) AS start_year, MAX(year) AS end_year").Joins("LEFT JOIN country_emissions ON country_emissions.country_id = countries.id").Group("countries.id").Find(&countries_api_response)
	c.JSON(200, gin.H{"data": countries_api_response})
}

func Country(c *gin.Context) {
	// Validate input
	var inputUri models.CountryEmissionApiRequestUri
	if err := c.ShouldBindUri(&inputUri); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var inputQuery models.CountryEmissionApiRequestQuery
	if err := c.ShouldBindQuery(&inputQuery); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var country_emission_api_response []models.CountryEmissionApiResponse

	tx := connections.Db.Table("country_emissions").Select("country_emissions.id, countries.id AS country_id, countries.name AS country_name, country_emissions.year, country_emissions.value, emission_categories.name AS emission_category").Joins("LEFT JOIN countries ON country_emissions.country_id = countries.id").Joins("JOIN emission_categories ON emission_categories.id = country_emissions.emission_category_id")

	// adds gas condition in query
	if inputQuery.Gas != "" {
		gases := strings.Split(inputQuery.Gas, ",")
		count := 0
		gas_where_condition := ""
		for _, gas := range gases {
			if count > 0 {
				gas_where_condition = gas_where_condition + " OR "
			}
			gas_condition := fmt.Sprintf("emission_categories.name LIKE '%v'", "%"+gas+"%")
			gas_where_condition = gas_where_condition + gas_condition
			count = count + 1
		}
		tx.Where("(" + gas_where_condition + ")")
	}

	// adds country condition in query
	tx.Where("countries.id=?", inputUri.CountryID)

	// adds start year condition in query
	if inputQuery.StartYear != 0 {
		tx.Where("year >= ?", inputQuery.StartYear)
	}

	// adds end year condition in query
	if inputQuery.EndYear != 0 {
		tx.Where("year <= ?", inputQuery.EndYear)
	}

	tx.Find(&country_emission_api_response)

	c.JSON(200, gin.H{"data": country_emission_api_response})
}

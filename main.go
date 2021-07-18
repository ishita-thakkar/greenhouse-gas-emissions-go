package main

import (
	"bluesky.com/greenhouse-gas-emissions/connections"
	"bluesky.com/greenhouse-gas-emissions/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func countries(c *gin.Context) {
	var countries_api_response []models.CountriesApiResponse
	connections.Db.Table("countries").Select("countries.id, countries.name, MIN(year) AS start_year, MAX(year) AS end_year").Joins("LEFT JOIN country_emissions ON country_emissions.country_id = countries.id").Group("countries.id").Find(&countries_api_response)
	c.JSON(200, gin.H{"data": countries_api_response})
}

func country(c *gin.Context) {
	country_id := c.Param("id")
	start_year_param := c.Query("startyear")
	end_year_param := c.Query("endyear")
	gas_param := c.Query("gas")

	var country_emission_api_response []models.CountryEmissionApiResponse

	tx := connections.Db.Table("country_emissions").Select("country_emissions.id, countries.id AS country_id, countries.name AS country_name, country_emissions.year, country_emissions.value").Joins("LEFT JOIN countries ON country_emissions.country_id = countries.id")

	if gas_param != "" {
		gases := strings.Split(gas_param, ",")
		tx.Joins("JOIN emission_categories ON emission_categories.id = country_emissions.emission_category_id")
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

	tx.Where("countries.id=?", country_id)

	if start_year_param != "" {
		tx.Where("year >= ?", start_year_param)
	}
	if end_year_param != "" {
		tx.Where("year <= ?", end_year_param)
	}

	tx.Find(&country_emission_api_response)

	c.JSON(200, gin.H{"data": country_emission_api_response})
}

func main() {
	router := gin.Default()
	router.GET("/ping", pingHandler)
	router.GET("/countries", countries)
	router.GET("/country/:id", country)
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

package models

// for reading csv data
type GreenhouseGas struct {
	ID            int
	CountryOrArea string
	Year          int
	Value         float64
	Category      string
}

type Country struct {
	ID   int
	Name string
}

type EmissionCategory struct {
	ID   int
	Name string
}

type CountryEmission struct {
	ID                 int
	CountryID          int
	Year               int `gorm:"index"`
	Value              float64
	EmissionCategoryID int
}

type CountriesApiResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartYear int    `json:"start_year"`
	EndYear   int    `json:"end_year"`
}

type CountryEmissionApiResponse struct {
	ID               int     `json:"id"`
	CountryID        int     `json:"country_id"`
	CountryName      string  `json:"country_name"`
	Year             int     `json:"year"`
	Value            float64 `json:"value"`
	EmissionCategory string  `json:"emission_category"`
}

type CountryEmissionApiRequestUri struct {
	CountryID int `uri:"id" binding:"required"`
}

type CountryEmissionApiRequestQuery struct {
	StartYear int    `form:"startyear" binding:"omitempty,min=0,max=9999"`
	EndYear   int    `form:"endyear" binding:"omitempty,min=0,max=9999,gtecsfield=StartYear"`
	Gas       string `form:"gas" binding:"omitempty,alphanum"`
}

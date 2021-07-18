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

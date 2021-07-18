package connections

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Db *gorm.DB

func init() {
	fmt.Println("Initialising database connections")
	var err error

	Db, err = gorm.Open(mysql.Open(os.Getenv("GHG_DB")), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}

	sqlDB, err := Db.DB()
	if err != nil {
		fmt.Println(err)
	}
	sqlDB.SetMaxOpenConns(100)
}

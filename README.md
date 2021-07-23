# Greenhouse gas emissions go
A script to format and insert data from an input csv to MySQL <br>
REST APIs to fetch data from MySQL with model validations for input parameters. <br>
Goroutines to concurrently insert part of the input data and channels for synchronizing the goroutines. <br>
GORM ORM and Gin web framework. <br>

## For local set up -
Set the environment variables<br>
Run populate_data.go to insert csv data in database <br>
Start server by running main <br>
Use /ping API to check for connectivity <br>
Remember to gofmt -w filename for formatting code before pushing

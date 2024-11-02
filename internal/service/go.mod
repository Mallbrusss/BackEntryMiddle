module service

go 1.23.1

replace internal/models => ./../../models

replace internal/repository => ./../repository

require internal/models v0.0.0-00010101000000-000000000000

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.14.0 // indirect
	gorm.io/gorm v1.25.12 // indirect
	internal/repository v0.0.0-00010101000000-000000000000
)

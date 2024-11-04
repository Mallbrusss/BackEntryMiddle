module service

go 1.23.1

replace internal/models => ./../../models

replace internal/repository => ./../repository

replace pkg/utils => ./../../pkg/utils

require internal/models v0.0.0-00010101000000-000000000000

require (
	github.com/google/uuid v1.6.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.28.0
	golang.org/x/text v0.19.0 // indirect
	gorm.io/gorm v1.25.12 // indirect
	internal/repository v0.0.0-00010101000000-000000000000
)

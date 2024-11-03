module postgres

go 1.23.1

replace internal/service => ./../service

replace pkg/utils => ./../../pkg/utils

replace internal/repository => ./../repository

replace internal/middleware => ./../middleware

replace internal/models => ./../../../models

replace backEntryMiddle/config => ./../../../config

replace internal/handlers => ./../handlers

replace backEntryMiddle/envloader => ./../../../envloader

require (
	backEntryMiddle/envloader v0.0.0-00010101000000-000000000000
	gorm.io/gorm v1.25.12
	internal/models v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.14.0 // indirect
	gorm.io/driver/postgres v1.5.9
)

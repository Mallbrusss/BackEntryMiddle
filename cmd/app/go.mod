module github.com/Mallbrusss/BackEntryMiddle/cmd/main

go 1.23.1

replace internal/server => ./../../internal/server

replace internal/service => ./../../internal/service

replace pkg/utils => ./../../pkg/utils

replace internal/repository => ./../../internal/repository

replace internal/middleware => ./../../internal/middleware

replace internal/models => ./../../models

replace backEntryMiddle/config => ./../../config

replace internal/handlers => ./../../internal/handlers

replace backEntryMiddle/envloader => ./../../envloader

replace internal/storage/postgres => ./../../internal/storage/postgres

require internal/server v0.0.0-00010101000000-000000000000

require (
	backEntryMiddle/envloader v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo/v4 v4.12.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	gorm.io/driver/postgres v1.5.9 // indirect
	gorm.io/gorm v1.25.12 // indirect
	internal/handlers v0.0.0-00010101000000-000000000000 // indirect
	internal/models v0.0.0-00010101000000-000000000000 // indirect
	internal/repository v0.0.0-00010101000000-000000000000 // indirect
	internal/service v0.0.0-00010101000000-000000000000 // indirect
	internal/storage/postgres v0.0.0-00010101000000-000000000000 // indirect
)

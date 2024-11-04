module handlers

go 1.23.1

replace internal/service => ./../service

replace pkg/utils => ./../../pkg/utils

replace internal/repository => ./../repository

replace internal/models => ./../../models

replace internal/middleware => ./../middleware

require (
	github.com/labstack/echo/v4 v4.12.0
	internal/models v0.0.0-00010101000000-000000000000
	internal/service v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	gorm.io/gorm v1.25.12 // indirect
	internal/repository v0.0.0-00010101000000-000000000000 // indirect
)

module github.com/jo-tbhac/kanban-api

go 1.13

replace local.packages/config => ./config

replace local.packages/db => ./db

replace local.packages/entity => ./entity

replace local.packages/handler => ./handler

replace local.packages/migration => ./migration

replace local.packages/repository => ./repository

replace local.packages/validator => ./validator

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0
	github.com/golang/protobuf v1.4.1 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/spf13/viper v1.6.3
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79
	golang.org/x/sys v0.0.0-20200501145240-bc7a7d42d5c3 // indirect
	local.packages/config v0.0.0-00010101000000-000000000000
	local.packages/db v0.0.0-00010101000000-000000000000
	local.packages/entity v0.0.0-00010101000000-000000000000 // indirect
	local.packages/handler v0.0.0-00010101000000-000000000000
	local.packages/migration v0.0.0-00010101000000-000000000000
	local.packages/repository v0.0.0-00010101000000-000000000000
	local.packages/validator v0.0.0-00010101000000-000000000000 // indirect
)

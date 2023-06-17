#!/bin/sh
go run migrate/migrate.go
go run seed/seed.go
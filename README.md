# e-Library REST API

A simple e-book loan management system built with Go.

## Features
- **Book Discovery**: Real-time availability checks.
- **Loan Management**: Borrow, Extend, and Return logic.
- **Observability**: Structured JSON logging to `stdout`.
- **Modular Design**: Interface-driven repository (Memory/Postgres).

## Tech Stack
- **Framework**: Gin Gonic
- **Logging**: zerolog
- **Persistence**: Postgres Support

## Quick Start
1. `go mod tidy`
2. `go run main.go`
3. `go test ./...`

## 📖 API Documentation
- **GET /Book?title=...** : Detail & Availability
- **POST /Borrow** : {name, title} -> 28-day loan
- **POST /Extend** : {name, title} -> +21 extraDays
- **POST /Return** : {name, title} -> Success

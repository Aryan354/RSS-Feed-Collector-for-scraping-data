# World News RSS Aggregator Server

## Overview
This project is an RSS server built in Go that aggregates world news from the BBC's RSS feed. It uses 10 concurrent goroutines to scrape data from the BBC's XML file available at [BBC RSS Feed](https://feeds.bbci.co.uk/news/world/rss.xml). The server uses Goose as a database migration tool and the SQLC library for querying. Unlike conventional ORM architectures, this project directly uses SQL queries. The relational database used is PostgreSQL.

Additionally, the project features API-based user authentication using API keys. Authentication is performed via the Authorization header with API keys generated at the time of storage.

## Features
- Concurrent scraping using 10 goroutines
- SQL-based querying with SQLC
- PostgreSQL for data storage
- Goose for database migration
- API-based authentication with API keys

## Prerequisites
- Go 1.19 or higher
- PostgreSQL
- Docker (for containerization)
- Git (for version control)

## Installation

### Install Go
Download and install Go from the [official Go website](https://golang.org/dl/).

### Install PostgreSQL
Follow the instructions on the [PostgreSQL website](https://www.postgresql.org/download/) to install PostgreSQL.

### Clone the Repository
```sh
git clone https://github.com/yourusername/your-repo.git
cd your-repo

### Configure the Database
Create a PostgreSQL database and update the connection string in your Go application configuration.

### Run Database Migrations
Install Goose for database migrations:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest

###Run the migrations:

```sh
Copy code
goose up
### Install Dependencies
### Install the required dependencies using Go modules:

````sh
Copy code
go mod tidy

###Running the Server
```sh
go run main.go

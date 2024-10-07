# URL Shortener API and Web Application

## Project Overview

This project provides a URL shortening service, which allows users to shorten URLs with an option for custom shortened links and expiration times. It is built with Go for both the API and the web interface and uses Redis as the primary database.

### Features
- Shorten URLs with optional customization and expiration times
- API built with Go (Golang) to handle URL shortening logic
- Web frontend created using HTML, CSS, and JavaScript for user interactions
- Redis is used to store and manage the shortened URLs
- Dockerized environment using Docker Compose for ease of deployment

## Tech Stack

- **Backend (API)**: Go
- **Backend (Web)**: Go
- **Frontend (Web)**: HTML, CSS, JavaScript
- **Database**: Redis
- **Containerization**: Docker, Docker Compose

## Project Structure

```bash
├── bin                 # Compiled binaries
├── cmd
│   ├── api             # API backend source code
│   │   └── Dockerfile.api
│   ├── web             # Web application source code
│   │   └── Dockerfile.web
├── data                # Data files (e.g., Redis storage)
├── docker-compose.yml  # Docker Compose configuration file
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
└── Makefile            # Build and run automation

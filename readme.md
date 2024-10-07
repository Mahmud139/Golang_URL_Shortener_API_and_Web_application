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
```

## Usage

### 1. Clone the repository:
```bash
gh repo clone Mahmud139/Golang_URL_Shortener_API_and_Web_application
cd your-repo
```

### 2. Build and run the application using Docker Compose:
```bash
docker-compose up --build
```
This will build and start both the API and web application in their respective Docker containers.

### 3. Access the web interface:
- Navigate to http://localhost:3005 in your browser to use the URL shortener.

### 4. API Endpoints:
- Navigate to http://localhost:3004 in your browser to see details of API Endpoints.
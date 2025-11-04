# Go Student API

![Go](https://img.shields.io/badge/Go-100%25-00ADD8?style=flat-square&logo=go)

A Go API project for managing student data with a clean architecture approach.

## Project Overview

This RESTful API service handles student management operations including:
- Student CRUD operations
- Course enrollment/management
- Student data validation
- Database persistence

## Project Structure

```
go-student-api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── configs/
│   ├── config.dev.yaml      # Development configuration
│   └── config.prod.yaml     # Production configuration
├── internal/
│   ├── db/
│   │   ├── config.go        # Database configuration
│   │   └── postgres.go      # Database connection management
│   ├── handlers/
│   │   ├── student.go       # Student HTTP handlers
│   │   └── course.go        # Course HTTP handlers
│   ├── models/
│   │   ├── student.go       # Student domain model
│   │   └── course.go        # Course domain model
│   ├── repository/
│   │   ├── postgres/        # PostgreSQL implementations
│   │   └── interfaces.go    # Repository interfaces
│   ├── service/
│   │   └── student.go       # Business logic layer
│   └── server/
│       └── server.go        # HTTP server setup
├── migrations/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── pkg/
│   └── validator/           # Shared validation utilities
├── scripts/
│   ├── migrations.sh        # Database migration helper
│   └── setup.sh            # Development setup script
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Go 1.19 or higher
- PostgreSQL
- Docker (optional)
- Make (optional)

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/manish-npx/go-student-api.git
cd go-student-api
```

2. Copy environment file and configure
```bash
cp .env.example .env
# Edit .env with your database credentials and other settings
```

3. Start the database (using Docker)
```bash
docker-compose up -d postgres
```

4. Run database migrations
```bash
./scripts/migrations.sh up
```

5. Build and run the application
```bash
go build -o bin/api ./cmd/api
./bin/api
```

## Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Server
PORT=8080
ENV=development
LOG_LEVEL=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=student_api
DB_SSLMODE=disable

# Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=300s
```

## API Endpoints

### Students
- `GET /api/v1/students` - List all students
- `GET /api/v1/students/{id}` - Get a specific student
- `POST /api/v1/students` - Create a new student
- `PUT /api/v1/students/{id}` - Update a student
- `DELETE /api/v1/students/{id}` - Delete a student

### Courses
- `GET /api/v1/courses` - List all courses
- `POST /api/v1/courses` - Create a new course
- `POST /api/v1/students/{id}/enroll` - Enroll student in courses

## Database Schema

### Students Table
```sql
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Courses Table
```sql
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Enrollments Table
```sql
CREATE TABLE enrollments (
    student_id INTEGER REFERENCES students(id),
    course_id INTEGER REFERENCES courses(id),
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (student_id, course_id)
);
```

## Development

### Running Tests
```bash
go test ./...
```

### Running with Docker
```bash
docker build -t student-api .
docker run -p 8080:8080 --env-file .env student-api
```

### Running Migrations
```bash
# Apply migrations
./scripts/migrations.sh up

# Rollback migrations
./scripts/migrations.sh down
```

## Project Layout Explanation

- `cmd/`: Contains the main application entry points
- `configs/`: Configuration files for different environments
- `internal/`: Private application code
  - `db/`: Database connection and configuration
  - `handlers/`: HTTP request handlers
  - `models/`: Domain models and business logic
  - `repository/`: Data access layer
  - `service/`: Business logic layer
  - `server/`: HTTP server setup and routing
- `migrations/`: Database migration files
- `pkg/`: Public libraries that can be used by external projects
- `scripts/`: Utility scripts for development and deployment

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

Manish - [@manish-npx](https://github.com/manish-npx)

Project Link: [https://github.com/manish-npx/go-student-api](https://github.com/manish-npx/go-student-api)
```

package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/manish-npx/go-student-api/internal/config"
	"github.com/manish-npx/go-student-api/internal/types"
)

type Postgres struct {
	DB *sql.DB
}

// -------------------------------------------------------------
// New() → Initializes the connection and ensures the students table
// -------------------------------------------------------------
func New(cfg config.Config) (*Postgres, error) {
	// ✅ Ensure database exists (auto-create if missing)
	if err := ensureDatabase(cfg); err != nil {
		return nil, fmt.Errorf("failed to ensure database: %w", err)
	}

	// ✅ Build DSN for the target DB
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	// ✅ Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("✅ Connected to PostgreSQL and ensured 'students' table")
	return &Postgres{DB: db}, nil
}

// -------------------------------------------------------------
// ensureDatabase() → Auto-creates DB if missing (when connected to postgres default DB)
// -------------------------------------------------------------
func ensureDatabase(cfg config.Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres system DB: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("CREATE DATABASE %s;", cfg.Postgres.DBName)
	_, err = db.Exec(query)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

// -------------------------------------------------------------
// CreateStudent() → Insert a student and return generated ID
// -------------------------------------------------------------
func (p *Postgres) CreateStudent(name, email string, age int) (int64, error) {
	var id int64
	err := p.DB.QueryRow(
		`INSERT INTO students (name, email, age)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		name, email, age,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to insert student: %w", err)
	}
	return id, nil
}

// -------------------------------------------------------------
// GetStudentById() → Fetch single student by ID
// -------------------------------------------------------------
func (p *Postgres) GetStudentById(id int64) (types.Student, error) {
	var student types.Student
	err := p.DB.QueryRow(
		`SELECT id, name, email, age
		 FROM students
		 WHERE id = $1`,
		id,
	).Scan(&student.ID, &student.Name, &student.Email, &student.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id: %d", id)
		}
		return types.Student{}, fmt.Errorf("failed to fetch student: %w", err)
	}

	return student, nil
}

// -------------------------------------------------------------
// GetStudents() → Fetch all students
// -------------------------------------------------------------
func (p *Postgres) GetStudents() ([]types.Student, error) {
	rows, err := p.DB.Query(`SELECT id, name, email, age FROM students ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query students: %w", err)
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return students, nil
}

// -------------------------------------------------------------
// UpdateStudentById() → Update student based on id
// -------------------------------------------------------------
func (p *Postgres) UpdateStudentById(id int64, name, email string, age int) (types.Student, error) {
	query := `UPDATE students SET name = $1, email = $2, age = $3 WHERE id = $4;`

	res, err := p.DB.Exec(query, name, email, age, id)
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to scan student: %w", err)
	}
	// Check if any rows were updated
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id: %d", id)
	}

	student, _ := p.GetStudentById(id)

	fmt.Println("Update student record is ", res)

	return student, nil
}

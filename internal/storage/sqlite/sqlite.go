package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/manish-npx/go-student-api/internal/config"
	"github.com/manish-npx/go-student-api/internal/types"
	_ "modernc.org/sqlite" // ✅ Pure-Go driver (no CGO)
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg config.Config) (*Sqlite, error) {
	if cfg.StoragePath == "" {
		return nil, fmt.Errorf("storage path not provided in config")
	}

	// ✅ Open or create SQLite DB file
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// ✅ Check connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	// ✅ Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("✅ SQLite connected and 'students' table ensured")

	return &Sqlite{Db: db}, nil
}

// -------------------------------------------------------------
// CreateStudent → Insert record
// -------------------------------------------------------------
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("prepare insert failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, fmt.Errorf("insert exec failed: %w", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to fetch last insert ID: %w", err)
	}

	return lastId, nil
}

// -------------------------------------------------------------
// GetStudentById → Fetch a single student by ID
// -------------------------------------------------------------
func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id: %d", id)
		}
		return types.Student{}, fmt.Errorf("query failed: %w", err)
	}

	return student, nil
}

// -------------------------------------------------------------
// GetStudents → Fetch all students
// -------------------------------------------------------------
func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return students, nil
}

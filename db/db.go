package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// User representa la estructura de la tabla users en la base de datos
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var db *sql.DB

// Connect establece la conexión con la base de datos PostgreSQL
func Connect() error {
	// Parámetros de conexión - estos deberían estar en variables de entorno en producción
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)

	// String de conexión
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Abrir la conexión
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Conexión exitosa a la base de datos")

	// Crear la tabla users si no existe
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		return err
	}

	return nil
}

// CreateUser inserta un nuevo usuario en la base de datos
func CreateUser(name, email, password string) (int, error) {
	var id int
	query := `INSERT INTO users (name, email, password) 
			  VALUES ($1, $2, $3) 
			  RETURNING id`

	err := db.QueryRow(query, name, email, password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetAllUsers retorna todos los usuarios de la base de datos
func GetAllUsers() ([]User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID busca un usuario por su ID
func GetUserByID(id int) (User, error) {
	var user User
	query := `SELECT id, name, email, created_at, updated_at 
			  FROM users 
			  WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

// UpdateUser actualiza la información de un usuario
func UpdateUser(id int, name, email string) error {
	query := `UPDATE users 
			  SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $3`

	_, err := db.Exec(query, name, email, id)
	return err
}

// DeleteUser elimina un usuario por su ID
func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Close cierra la conexión con la base de datos
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

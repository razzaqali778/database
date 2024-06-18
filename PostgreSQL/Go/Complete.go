package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

// Connect to the database
func connect() {
	var err error
	dsn := "postgres://yourusername:yourpassword@localhost:5432/yourdbname"
	pool, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Println("Connected to PostgreSQL database")
}

// Close the database connection
func close() {
	pool.Close()
}

// Utility function to execute queries
func executeQuery(query string, params ...interface{}) (pgx.Rows, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire connection: %v", err)
	}
	defer conn.Release()
	return conn.Query(context.Background(), query, params...)
}

// Schema Definition
func createTables() {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			age INT
		);
	`

	createOrdersTable := `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			product VARCHAR(100),
			amount INT
		);
	`

	_, err := executeQuery(createUsersTable)
	if err != nil {
		log.Fatalf("Unable to create users table: %v\n", err)
	}

	_, err = executeQuery(createOrdersTable)
	if err != nil {
		log.Fatalf("Unable to create orders table: %v\n", err)
	}

	log.Println("Tables created successfully")
}

// CRUD Operations
func createUser(name, email string, age int) {
	query := "INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id"
	rows, err := executeQuery(query, name, email, age)
	if err != nil {
		log.Fatalf("Unable to create user: %v\n", err)
	}
	defer rows.Close()

	var id int
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatalf("Unable to scan user ID: %v\n", err)
		}
		log.Printf("User created with ID: %d\n", id)
	}
}

func getUsers() {
	query := "SELECT id, name, email, age FROM users"
	rows, err := executeQuery(query)
	if err != nil {
		log.Fatalf("Unable to get users: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, age int
		var name, email string
		err := rows.Scan(&id, &name, &email, &age)
		if err != nil {
			log.Fatalf("Unable to scan user: %v\n", err)
		}
		log.Printf("User: ID=%d, Name=%s, Email=%s, Age=%d\n", id, name, email, age)
	}
}

func updateUser(id int, name, email string, age int) {
	query := "UPDATE users SET name = $1, email = $2, age = $3 WHERE id = $4 RETURNING id"
	rows, err := executeQuery(query, name, email, age, id)
	if err != nil {
		log.Fatalf("Unable to update user: %v\n", err)
	}
	defer rows.Close()

	if rows.Next() {
		var updatedID int
		err := rows.Scan(&updatedID)
		if err != nil {
			log.Fatalf("Unable to scan updated user ID: %v\n", err)
		}
		log.Printf("User updated with ID: %d\n", updatedID)
	}
}

func deleteUser(id int) {
	query := "DELETE FROM users WHERE id = $1 RETURNING id"
	rows, err := executeQuery(query, id)
	if err != nil {
		log.Fatalf("Unable to delete user: %v\n", err)
	}
	defer rows.Close()

	if rows.Next() {
		var deletedID int
		err := rows.Scan(&deletedID)
		if err != nil {
			log.Fatalf("Unable to scan deleted user ID: %v\n", err)
		}
		log.Printf("User deleted with ID: %d\n", deletedID)
	}
}

// Query Operators
func queryOperators() {
	query := `
		SELECT * FROM users
		WHERE age >= 18 AND age <= 30
			AND name IN ('Alice', 'Bob')
			AND (age < 25 OR name = 'Charlie')
			AND age > 20 AND name <> 'Dave'
	`
	rows, err := executeQuery(query)
	if err != nil {
		log.Fatalf("Unable to execute query operators: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, age int
		var name, email string
		err := rows.Scan(&id, &name, &email, &age)
		if err != nil {
			log.Fatalf("Unable to scan user: %v\n", err)
		}
		log.Printf("User: ID=%d, Name=%s, Email=%s, Age=%d\n", id, name, email, age)
	}
}

// Update Operators
func updateOperators() {
	query := `
		UPDATE users
		SET age = age + 1,
			name = 'Updated Name',
			age = CASE WHEN age > 30 THEN age - 1 ELSE age END,
			email = NULL
		WHERE name = 'Alice'
		RETURNING id
	`
	rows, err := executeQuery(query)
	if err != nil {
		log.Fatalf("Unable to execute update operators: %v\n", err)
	}
	defer rows.Close()

	if rows.Next() {
		var updatedID int
		err := rows.Scan(&updatedID)
		if err != nil {
			log.Fatalf("Unable to scan updated user ID: %v\n", err)
		}
		log.Printf("User updated with ID: %d\n", updatedID)
	}
}

// Aggregation Functions
func aggregationFunctions() {
	query := `
		SELECT user_id, SUM(amount) as total_amount, AVG(amount) as average_amount
		FROM orders
		GROUP BY user_id
		HAVING SUM(amount) > 100
		ORDER BY total_amount DESC
	`
	rows, err := executeQuery(query)
	if err != nil {
		log.Fatalf("Unable to execute aggregation functions: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		var totalAmount, averageAmount float64
		err := rows.Scan(&userID, &totalAmount, &averageAmount)
		if err != nil {
			log.Fatalf("Unable to scan aggregation result: %v\n", err)
		}
		log.Printf("Aggregation: UserID=%d, TotalAmount=%.2f, AverageAmount=%.2f\n", userID, totalAmount, averageAmount)
	}
}

// Joins
func getUsersWithOrders() {
	query := `
		SELECT u.name, u.email, o.product, o.amount
		FROM users u
		JOIN orders o ON u.id = o.user_id
	`
	rows, err := executeQuery(query)
	if err != nil {
		log.Fatalf("Unable to execute join query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, email, product string
		var amount int
		err := rows.Scan(&name, &email, &product, &amount)
		if err != nil {
			log.Fatalf("Unable to scan join result: %v\n", err)
		}
		log.Printf("Join: Name=%s, Email=%s, Product=%s, Amount=%d\n", name, email, product, amount)
	}
}

// Indexing
func createIndexes() {
	createIndex := `
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)
	`
	dropIndex := `
		DROP INDEX IF EXISTS idx_users_email
	`

	_, err := executeQuery(createIndex)
	if err != nil {
		log.Fatalf("Unable to create index: %v\n", err)
	}
	log.Println("Index created successfully")

	_, err = executeQuery(dropIndex)
	if err != nil {
		log.Fatalf("Unable to drop index: %v\n", err)
	}
	log.Println("Index dropped successfully")
}

// Transactions
func executeTransaction() {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire connection: %v\n", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v\n", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
			log.Fatalf("Transaction rolled back: %v\n", err)
		} else {
			tx.Commit(context.Background())
			log.Println("Transaction committed successfully")
		}
	}()

	_, err = tx.Exec(context.Background(), "INSERT INTO users (name, email, age) VALUES ($1, $2, $3)", "Charlie", "charlie@example.com", 22)
	if err != nil {
		return
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO users (name, email, age) VALUES ($1, $2, $3)", "Dana", "dana@example.com", 28)
	if err != nil {
		return
	}
}

// Miscellaneous Operations
func miscellaneousOperations() {
	explainQuery := "EXPLAIN SELECT * FROM users"
	rows, err := executeQuery(explainQuery)
	if err != nil {
		log.Fatalf("Unable to execute explain query: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var explanation string
		err := rows.Scan(&explanation)
		if err != nil {
			log.Fatalf("Unable to scan explain result: %v\n", err)
		}
		log.Println("Explain query result:", explanation)
	}

	vacuumQuery := "VACUUM"
	_, err = executeQuery(vacuumQuery)
	if err != nil {
		log.Fatalf("Unable to execute vacuum: %v\n", err)
	}
	log.Println("VACUUM executed successfully")

	analyzeQuery := "ANALYZE"
	_, err = executeQuery(analyzeQuery)
	if err != nil {
		log.Fatalf("Unable to execute analyze: %v\n", err)
	}
	log.Println("ANALYZE executed successfully")

	copyToFileQuery := "COPY users TO '/path/to/file.csv' WITH (FORMAT CSV)"
	_, err = executeQuery(copyToFileQuery)
	if err != nil {
		log.Fatalf("Unable to copy to file: %v\n", err)
	}
	log.Println("Data exported to CSV file successfully")
}

// Main function to run the examples
func main() {
	connect()
	defer close()

	// Create tables
	createTables()

	// Run CRUD operations
	createUser("Alice", "alice@example.com", 25)
	createUser("Bob", "bob@example.com", 30)
	getUsers()
	updateUser(1, "Alice Smith", "alice.smith@example.com", 26)
	deleteUser(2)
	getUsers()

	// Run Query Operators
	queryOperators()

	// Run Update Operators
	updateOperators()

	// Run Aggregation Functions
	aggregationFunctions()

	// Run Joins
	getUsersWithOrders()

	// Run Indexing
	createIndexes()

	// Run Transactions
	executeTransaction()

	// Run Miscellaneous Operations
	miscellaneousOperations()
}

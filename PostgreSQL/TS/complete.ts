import { Pool } from 'pg'

// Database connection configuration
const pool = new Pool({
  user: 'yourusername',
  host: 'localhost',
  database: 'yourdbname',
  password: 'yourpassword',
  port: 5432,
})

pool.on('connect', () => {
  console.log('Connected to the PostgreSQL database')
})

pool.on('error', (err) => {
  console.error('Unexpected error on idle client', err)
  process.exit(-1)
})

// Utility function to execute queries
const executeQuery = async (query: string, params?: any[]) => {
  const client = await pool.connect()
  try {
    const res = await client.query(query, params)
    return res
  } catch (err) {
    console.error('Error executing query:', err)
    throw err
  } finally {
    client.release()
  }
}

// Schema Definition
const createTables = async () => {
  const createUsersTable = `
    CREATE TABLE IF NOT EXISTS users (
      id SERIAL PRIMARY KEY,
      name VARCHAR(100) NOT NULL,
      email VARCHAR(100) UNIQUE NOT NULL,
      age INT
    );
  `

  const createOrdersTable = `
    CREATE TABLE IF NOT EXISTS orders (
      id SERIAL PRIMARY KEY,
      user_id INT REFERENCES users(id),
      product VARCHAR(100),
      amount INT
    );
  `

  await executeQuery(createUsersTable)
  await executeQuery(createOrdersTable)
  console.log('Tables created successfully')
}

// CRUD Operations
const createUser = async (name: string, email: string, age: number) => {
  const query =
    'INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING *'
  const res = await executeQuery(query, [name, email, age])
  console.log('User created:', res.rows[0])
}

const getUsers = async () => {
  const query = 'SELECT * FROM users'
  const res = await executeQuery(query)
  console.log('Users:', res.rows)
}

const updateUser = async (
  id: number,
  name: string,
  email: string,
  age: number
) => {
  const query =
    'UPDATE users SET name = $1, email = $2, age = $3 WHERE id = $4 RETURNING *'
  const res = await executeQuery(query, [name, email, age, id])
  console.log('User updated:', res.rows[0])
}

const deleteUser = async (id: number) => {
  const query = 'DELETE FROM users WHERE id = $1 RETURNING *'
  const res = await executeQuery(query, [id])
  console.log('User deleted:', res.rows[0])
}

// Query Operators
const queryOperators = async () => {
  const query = `
    SELECT * FROM users
    WHERE age >= 18 AND age <= 30
      AND name IN ('Alice', 'Bob')
      AND (age < 25 OR name = 'Charlie')
      AND age > 20 AND name <> 'Dave'
      AND reactions IS NOT NULL;
  `
  const res = await executeQuery(query)
  console.log('Query operator results:', res.rows)
}

// Update Operations
const updateOperators = async () => {
  const query = `
    UPDATE users
    SET age = age + 1,
        name = 'Updated Name',
        age = CASE WHEN age > 30 THEN age - 1 ELSE age END,
        reactions = jsonb_set(reactions, '{key}', '"value"'),
        email = NULL
    WHERE name = 'Alice'
    RETURNING *;
  `
  const res = await executeQuery(query)
  console.log('Update operators result:', res.rows)
}

// Aggregation Functions
const aggregationFunctions = async () => {
  const query = `
    SELECT user_id, SUM(amount) as total_amount, AVG(amount) as average_amount
    FROM orders
    GROUP BY user_id
    HAVING SUM(amount) > 100
    ORDER BY total_amount DESC;
  `
  const res = await executeQuery(query)
  console.log('Aggregation results:', res.rows)
}

// Joins
const getUsersWithOrders = async () => {
  const query = `
    SELECT u.name, u.email, o.product, o.amount
    FROM users u
    JOIN orders o ON u.id = o.user_id;
  `
  const res = await executeQuery(query)
  console.log('Users with orders:', res.rows)
}

// Indexing
const createIndexes = async () => {
  const createIndex =
    'CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)'
  const dropIndex = 'DROP INDEX IF EXISTS idx_users_email'

  await executeQuery(createIndex)
  console.log('Index created successfully')

  await executeQuery(dropIndex)
  console.log('Index dropped successfully')
}

// Transactions
const executeTransaction = async () => {
  const client = await pool.connect()
  try {
    await client.query('BEGIN')
    await client.query(
      'INSERT INTO users (name, email, age) VALUES ($1, $2, $3)',
      ['Charlie', 'charlie@example.com', 22]
    )
    await client.query(
      'INSERT INTO users (name, email, age) VALUES ($1, $2, $3)',
      ['Dana', 'dana@example.com', 28]
    )
    await client.query('COMMIT')
    console.log('Transaction completed successfully')
  } catch (err) {
    await client.query('ROLLBACK')
    console.error('Error executing transaction, rolled back', err)
  } finally {
    client.release()
  }
}

// Miscellaneous Operations
const miscellaneousOperations = async () => {
  const explainQuery = 'EXPLAIN SELECT * FROM users'
  const explainResult = await executeQuery(explainQuery)
  console.log('Explain query result:', explainResult.rows)

  const vacuumQuery = 'VACUUM'
  await executeQuery(vacuumQuery)
  console.log('VACUUM executed successfully')

  const analyzeQuery = 'ANALYZE'
  await executeQuery(analyzeQuery)
  console.log('ANALYZE executed successfully')

  const copyToFileQuery = "COPY users TO '/path/to/file.csv' WITH (FORMAT CSV)"
  await executeQuery(copyToFileQuery)
  console.log('Data exported to CSV file successfully')
}

// Main function to run the examples
;(async () => {
  await createTables()

  // Run CRUD operations
  await createUser('Alice', 'alice@example.com', 25)
  await createUser('Bob', 'bob@example.com', 30)
  await getUsers()
  await updateUser(1, 'Alice Smith', 'alice.smith@example.com', 26)
  await deleteUser(2)
  await getUsers()

  // Run Query Operators
  await queryOperators()

  // Run Update Operators
  await updateOperators()

  // Run Aggregation Functions
  await aggregationFunctions()

  // Run Joins
  await getUsersWithOrders()

  // Run Indexing
  await createIndexes()

  // Run Transactions
  await executeTransaction()

  // Run Miscellaneous Operations
  await miscellaneousOperations()
})()

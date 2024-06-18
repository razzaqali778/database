# PostgreSQL CRUD Operations and Advanced Features

This guide provides an overview of CRUD operations, query operators, update operations, and other useful PostgreSQL features.

## CRUD Operations

### Create

- **INSERT INTO table_name (columns) VALUES (values);**: Inserts a single row into a table.
- **COPY table_name FROM 'file_path' WITH (FORMAT CSV);**: Inserts multiple rows from a CSV file into a table.

### Read

- **SELECT columns FROM table_name WHERE condition;**: Queries for rows in a table.
- **SELECT COUNT(\*) FROM table_name WHERE condition;**: Counts the number of rows in a table.
- **SELECT DISTINCT column FROM table_name;**: Selects distinct values from a column.

### Update

- **UPDATE table_name SET column = value WHERE condition;**: Updates rows in a table.
- **UPDATE table_name SET column = value FROM other_table WHERE table_name.id = other_table.id;**: Updates rows with values from another table.

### Delete

- **DELETE FROM table_name WHERE condition;**: Deletes rows from a table.
- **TRUNCATE table_name;**: Deletes all rows from a table.

## Query Operators

### Comparison

- **=**: Matches values that are equal to a specified value.
- **<>** or **!=**: Matches all values that are not equal to a specified value.
- **>**: Matches values that are greater than a specified value.
- **>=**: Matches values that are greater than or equal to a specified value.
- **<**: Matches values that are less than a specified value.
- **<=**: Matches values that are less than or equal to a specified value.
- **IN (list)**: Matches any of the values specified in a list.
- **NOT IN (list)**: Matches none of the values specified in a list.

### Logical

- **AND**: Joins query clauses with a logical AND returns all rows that match the conditions of both clauses.
- **OR**: Joins query clauses with a logical OR returns all rows that match the conditions of either clause.
- **NOT**: Inverts the effect of a query expression and returns rows that do not match the query expression.

### Pattern Matching

- **LIKE**: Matches values that match a specified pattern.
- **ILIKE**: Case-insensitive match to a specified pattern.
- **SIMILAR TO**: Matches values against a SQL regular expression.

### Other

- **IS NULL**: Matches rows where the specified column is null.
- **IS NOT NULL**: Matches rows where the specified column is not null.
- **BETWEEN**: Matches values within a specified range.
- **EXISTS**: Matches rows where a subquery returns one or more rows.

## Update Operations

### Field Update

- **SET column = value**: Sets the value of a column in a table.
- **SET column = value + amount**: Increments the value of a column by the specified amount.
- **SET column = value \* factor**: Multiplies the value of a column by the specified factor.
- **SET column = NULL**: Sets the value of a column to NULL.

### JSON Update

- **column->'key'**: Retrieves JSON object field by key.
- **column->>'key'**: Retrieves JSON object field as text by key.
- **column#>>'{path}'**: Retrieves JSON object field by path elements.
- **column || json_object**: Merges JSON objects.
- **column - 'key'**: Deletes a key from JSON object.

## Aggregation Functions

- **SUM(column)**: Returns the sum of the values.
- **AVG(column)**: Returns the average of the values.
- **MIN(column)**: Returns the minimum value.
- **MAX(column)**: Returns the maximum value.
- **COUNT(column)**: Returns the count of non-null values.
- **GROUP BY column**: Groups rows that have the same values into summary rows.
- **HAVING condition**: Filters groups based on a condition.

## Joins

- **INNER JOIN**: Returns rows when there is a match in both tables.
- **LEFT JOIN**: Returns all rows from the left table, and the matched rows from the right table.
- **RIGHT JOIN**: Returns all rows from the right table, and the matched rows from the left table.
- **FULL OUTER JOIN**: Returns rows when there is a match in one of the tables.

## Indexing

- **CREATE INDEX index_name ON table_name (columns);**: Creates an index on a table.
- **CREATE UNIQUE INDEX index_name ON table_name (columns);**: Creates a unique index on a table.
- **DROP INDEX index_name;**: Removes an index from a table.
- **REINDEX TABLE table_name;**: Rebuilds all indexes on a table.

## Transactions

- **BEGIN;**: Starts a transaction block.
- **COMMIT;**: Commits the current transaction.
- **ROLLBACK;**: Rolls back the current transaction.

## Miscellaneous

- **EXPLAIN query;**: Shows the execution plan of a query.
- **VACUUM;**: Reclaims storage occupied by dead tuples.
- **ANALYZE;**: Collects statistics about the contents of tables in the database.
- **COPY table_name TO 'file_path' WITH (FORMAT CSV);**: Exports table data to a CSV file.

## Important Considerations

- **Indexes**: Use indexes to optimize query performance.
- **Transactions**: Use transactions to ensure data consistency.
- **Schema Design**: Consider your application’s query patterns when designing your schema.
- **Error Handling**: Always implement error handling for database operations.
- **Security**: Implement proper access controls and encryption.

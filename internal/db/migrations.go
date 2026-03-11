package db

import "fmt"

const (
	createTableAccountsDDL = `CREATE TABLE IF NOT EXISTS accounts
(
	id SERIAL PRIMARY KEY,
	balance FLOAT,
	owner VARCHAR
);`

	createTableTransactionsDDL = `CREATE TABLE IF NOT EXISTS transactions
(
	id SERIAL PRIMARY KEY,
	amount FLOAT,
	type VARCHAR,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	account_id INT REFERENCES accounts (id)
);`
)

func RunMigrations() error {
	_, err := db.Exec(createTableAccountsDDL)
	if err != nil {
		fmt.Println("error creating accounts table")
		return err
	}
	
	_, err = db.Exec(createTableTransactionsDDL)
	if err != nil {
		fmt.Println("error creating transactions table")
		return err
	}

	return nil
}

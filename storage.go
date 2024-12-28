package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	createAccount(*Account) error
	deleteAccount(int) error
	updateAccount(*Account) error
	getAccountById(string) (*Account, error)
	getAllAccounts() ([]*Account, error)
}

type PostgressStore struct {
	db *sql.DB
}

func (s *PostgressStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgressStore) createAccountTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		number INT,
		balance NUMERIC,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`)
	return err
}

func newPostgressStore() (*PostgressStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &PostgressStore{
		db: db,
	}, nil
}

func (s *PostgressStore) createAccount(acc *Account) error {
	query := `
	insert into accounts (first_name,last_name,number,balance,created_at)
	values($1,$2,$3,$4,$5)`
	response, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", response)
	return nil
}
func (s *PostgressStore) updateAccount(*Account) error {
	return nil
}
func (s *PostgressStore) deleteAccount(id int) error {
	return nil
}
func (s *PostgressStore) getAccountById(id string) (*Account, error) {
	row := s.db.QueryRow("SELECT id, first_name, last_name, number, balance, created_at FROM accounts WHERE id=$1", id)

	var account Account

	err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with id %s not found", id)
		}
		return nil, err
	}

	return &account, nil
}

func (s *PostgressStore) getAllAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT id, first_name, last_name, number, balance, created_at FROM accounts;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []*Account

	for rows.Next() {
		var acc Account
		err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

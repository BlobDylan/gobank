package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type storage interface {
	CreateAccount(*account) error
	GetAccount(id int) (*account, error)
	GetAccounts() ([]*account, error)
	DeleteAccount(id int) error
	Transfer(from, to int, amount int64) error
}

type postgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*postgresStorage, error) {
	connstr := "user=postgres dbname=postgres password=1234 sslmode=disable"
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &postgresStorage{db: db}, nil
}
func (s *postgresStorage) Init() error {
	return s.CreateAccountTable()
}

func (s *postgresStorage) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		number BIGINT,
		email TEXT,
		balance BIGINT,
		created_at TIMESTAMP DEFAULT NOW()
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *postgresStorage) CreateAccount(acc *account) error {
	query := `INSERT INTO accounts (number, email, balance, created_at) VALUES ($1, $2, $3, $4)`

	resp, err := s.db.Query(query, acc.Number, acc.Email, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *postgresStorage) GetAccount(id int) (*account, error) {
	return nil, nil
}

func (s *postgresStorage) DeleteAccount(id int) error {
	return nil
}

func (s *postgresStorage) Transfer(from, to int, amount int64) error {
	return nil
}

func (s *postgresStorage) GetAccounts() ([]*account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts")

	if err != nil {
		return nil, err
	}

	accounts := []*account{}

	for rows.Next() {
		acc := &account{}
		if err := rows.Scan(&acc.ID, &acc.Number, &acc.Email, &acc.Balance, &acc.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

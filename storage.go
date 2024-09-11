package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type storage interface {
	CreateAccount(*account) error
	GetAccount(id int) (*account, error)
	GetAccountByEmail(email string) (*account, error)
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
		encrypted_pwd TEXT,
		balance BIGINT,
		created_at TIMESTAMP DEFAULT NOW()
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *postgresStorage) CreateAccount(acc *account) error {
	query := `INSERT INTO accounts (number, email, encrypted_pwd, balance, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(query, acc.Number, acc.Email, acc.EncryptedPwd, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *postgresStorage) GetAccount(id int) (*account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		return nil, err

	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *postgresStorage) GetAccountByEmail(Email string) (*account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE email = $1", Email)
	if err != nil {
		return nil, err

	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %s not found", Email)
}

func (s *postgresStorage) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id = $1", id)
	return err
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
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*account, error) {
	acc := new(account)
	err := rows.Scan(&acc.ID, &acc.Number, &acc.Email, &acc.EncryptedPwd, &acc.Balance, &acc.CreatedAt)

	return acc, err
}

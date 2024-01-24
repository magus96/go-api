package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(id int) error
	UpdateAccount() error
	GetAccounts() ([]*Account, error)
	GetAccountbyID(id int) (*Account, error)
	GetAccountbyNumber(num int64) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func newpostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=mybankapi sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account(
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance numeric,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `
	select * from account`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(&account.ID, &account.Firstname, &account.Lastname, &account.Number, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account
	(first_name, last_name, number, balance, created_at) 
	values($1, $2, $3, $4, $5)`
	_, err := s.db.Query(query, acc.Firstname, acc.Lastname, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `delete from account where id=$1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) UpdateAccount() error {
	return nil
}

func (s *PostgresStore) GetAccountbyNumber(num int64) (*Account, error) {
	query := "select * from account where number=$1"
	rows, err := s.db.Query(query, num)
	if err != nil {
		return nil, err
	}
	account := new((Account))
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Firstname, &account.Lastname, &account.Number, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Account with number %d not found", num)
		}
	}
	return account, nil
}

func (s *PostgresStore) GetAccountbyID(id int) (*Account, error) {
	query := `select * from account where id=$1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	account := new(Account)
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Firstname, &account.Lastname, &account.Number, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Account %d not found", id)
		}
	}
	return account, nil
}

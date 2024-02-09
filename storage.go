package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	createAccount(*Account) error
	deleteAccount(int) error
	updateAccount(*Account) error
	getAccounts() ([]*Account, error)
	getAccountByID(int) (*Account, error)
}

type postGresStorage struct {
	db *sql.DB
}

func NewPostGresStorage() (*postGresStorage, error) {

	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postGresStorage{db: db}, nil
}

func (s *postGresStorage) init() error {
	return s.createAccountTable()

}

func (s *postGresStorage) createAccountTable() error {
    query := `CREATE TABLE IF NOT EXISTS accounts (
        id SERIAL PRIMARY KEY,
        first_name varchar(255),
        last_name varchar(255),
        age int,
        balance serial,
        bank_number serial,
        createdAt timestamp
    )`
    _, err := s.db.Exec(query)
    return err
}


func (s *postGresStorage) createAccount(a *Account) error {

	query := (`
		INSERT INTO accounts
		(first_name, last_name,age, balance, bank_number, createdAt)
		VALUES ($1, $2, $3, $4, $5, $6)
	`)

	resp, err := s.db.Query(query, a.FirstName, a.LastName, a.Age, a.Balance, a.BankNumber, a.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Print(resp)

	return nil
}

func (s *postGresStorage) deleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id=$1", id)
	if err != nil {	
	return err	}

	return nil
}

func (s *postGresStorage) updateAccount(a *Account) error {
	return nil
}

func (s *postGresStorage) getAccounts() ([]*Account, error) {

	rows,err:=s.db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}

	accounts:=[]*Account{}
	for rows.Next(){
		account,err:=scanIntoAccounts(rows)
		if err != nil {return nil, err}
		accounts=append(accounts, account)
	}


	return accounts, nil
}


func (s *postGresStorage) getAccountByID(id int) (*Account, error) {
	rows,err:=s.db.Query("SELECT * FROM accounts WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next(){
		return scanIntoAccounts(rows)
	}
	return nil, err
}

func scanIntoAccounts(rows *sql.Rows) (*Account,error){
	account:=&Account{}

	err:=rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Age, &account.Balance, &account.BankNumber, &account.CreatedAt)	

	return account, err
}

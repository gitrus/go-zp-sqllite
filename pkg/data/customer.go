package data

import (
	"database/sql"

	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
)

type CustomerDBClient struct {
	db *sql.DB
}

func NewCustomerDBClient(db *sql.DB) *CustomerDBClient {
	return &CustomerDBClient{db: db}
}

func (cdb *CustomerDBClient) Get(id int) (e.Customer, error) {
	cl := e.Customer{}

	row := cdb.db.QueryRow("SELECT id, fio, login, birthday, email FROM clients WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&cl.ID, &cl.FIO, &cl.Login, &cl.Birthday, &cl.Email)
	if err != nil {
		return cl, err
	}

	return cl, nil
}

func (cdb *CustomerDBClient) Create(client e.Customer) (int, error) {
	res, err := cdb.db.Exec("INSERT INTO clients (fio, login, birthday, email) VALUES (:fio, :login, :birthday, :email)",
		sql.Named("fio", client.FIO),
		sql.Named("login", client.Login),
		sql.Named("birthday", client.Birthday),
		sql.Named("email", client.Email))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (cdb *CustomerDBClient) Delete(id int) error {
	_, err := cdb.db.Exec("DELETE FROM clients WHERE id = :id", sql.Named("id", id))

	return err
}

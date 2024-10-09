package data_test

import (
	"testing"

	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	// настройте подключение к БД

	clientID := 1

	// напиши тест здесь
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	// настройте подключение к БД

	clientID := -1

	// напиши тест здесь
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	// настройте подключение к БД

	cl := e.Customer{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// напиши тест здесь
}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	// настройте подключение к БД

	cl := e.Customer{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// напиши тест здесь
}

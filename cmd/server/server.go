package main

import (
	"comsrv/pkg/api"
	"comsrv/pkg/storage"
	"comsrv/pkg/storage/pgdb"
	"log"
	"net/http"
	"os"
)

// Сервер newssrv.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера
	srv := server{}

	//  Создаём объект базы данных PostgreSQL.
	pwd := os.Getenv("pgpass")
	connstr := "postgres://postgres:" + pwd + "@0.0.0.0/comsrv"
	db, err := pgdb.New(connstr)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера БД
	srv.db = db

	// Освобождаем ресурс
	defer srv.db.Close()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8082 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов.
	log.Println("HTTP server is started on localhost:8082")
	defer log.Println("HTTP server has been stopped")
	log.Fatal(http.ListenAndServe(":8082", srv.api.Router()))
}

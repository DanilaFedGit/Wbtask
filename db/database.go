package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// структура для функции ConnectDB, в которой хранится необходимая информация для подключения
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// структура для получения данных из бд
type Check struct {
	Check_id   string `db:"check_id"`
	Check_data string `db:"check_data"`
}

// функция для подключения к бд
func ConnectDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode =%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// функция для создания словаря(кеша)
func GetData(db *sqlx.DB) (map[string][]byte, error) {
	sliceofcheck := make([]Check, 0, 0)
	err := db.Select(&sliceofcheck, `select * from wb_check`)
	mapofcheck := make(map[string][]byte)
	for _, i := range sliceofcheck {
		mapofcheck[i.Check_id] = []byte(i.Check_data)
	}
	return mapofcheck, err

}

// функция для вставки данных в бд
func InsertData(db *sqlx.DB, id string, json_data []byte) {
	tx := db.MustBegin()
	tx.MustExec(`INSERT into wb_check (check_id,check_data) VALUES ($1,$2)`, id, json_data)
	tx.Commit()
}

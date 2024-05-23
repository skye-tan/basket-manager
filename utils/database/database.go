package database

import (
	"errors"
	"time"

	"github.com/skye-tan/basket-manager/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Basket struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Data      string    `json:"data"`
	State     string    `json:"state"`
}

var db *gorm.DB

func InitializeDatabase() error {
	var err error

	dsn := "host=localhost user=postgres password=a dbname=web port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return err
}

func GetBaskets() []Basket {
	var baskets []Basket
	db.Raw("SELECT * FROM baskets;").Scan(&baskets)
	return baskets
}

func CreateBasket(data string, state string) (uint, error) {
	if state != "COMPLETED" && state != "PENDING" {
		return 0, errors.New(custom_error.INVALID_STATE)
	}

	var id uint
	db.Raw("INSERT INTO baskets (created_at, updated_at, data, state) VALUES(?, ?, ?, ?) RETURNING id;",
		time.Now(), time.Now(), data, state).Scan(&id)
	return id, nil
}

func UpdateBasket(id uint, data string, state string) error {
	if state != "COMPLETED" && state != "PENDING" {
		return errors.New(custom_error.INVALID_STATE)
	}

	db.Exec("UPDATE baskets SET updated_at = ?, data = ?, state = ? WHERE id = ?;",
		time.Now(), data, state, id)
	return nil
}

func GetBasket(id uint) (Basket, error) {
	var basket Basket
	tx := db.Raw("SELECT * FROM baskets WHERE id = ?;", id).Scan(&basket)
	if tx.RowsAffected != 1 {
		return basket, errors.New(custom_error.INVALID_ID)
	}
	return basket, nil
}

func DeleteBasket(id uint) error {
	tx := db.Exec("DELETE FROM baskets WHERE id = ?;", id)
	if tx.RowsAffected != 1 {
		return errors.New(custom_error.INVALID_ID)
	}
	return nil
}

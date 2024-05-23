package database

import (
	"errors"
	"time"

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
	db.Raw("SELECT id, created_at, updated_at, data, state FROM baskets;").Scan(&baskets)
	return baskets
}

func CreateBasket(data string, state string) uint {
	var id uint
	db.Raw("INSERT INTO baskets (created_at, updated_at, data, state) VALUES(?, ?, ?, ?) RETURNING id;",
		time.Now(), time.Now(), data, state).Scan(&id)
	return id
}

func UpdateBasket(id uint, data string, state string) {
	db.Exec("UPDATE baskets SET updated_at = ?, data = ?, state = ? WHERE id = ?;",
		time.Now(), data, state, id)
}

func GetBasket(id uint) (Basket, error) {
	var basket Basket
	tx := db.Raw("SELECT * FROM baskets WHERE id = ?;", id).Scan(&basket)
	if tx.RowsAffected == 0 {
		return basket, errors.New("invalid id")
	}
	return basket, nil
}

func DeleteBasket(id uint) {
	db.Exec("DELETE FROM baskets WHERE id = ?;", id)
}

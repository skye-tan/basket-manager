package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	custom_error "github.com/skye-tan/basket-manager/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	PENDING   = "PENDING"
	COMPLETED = "COMPLETED"
)

type Basket struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Data      JSONB     `json:"data"`
	State     string    `json:"state"`
}

type JSONB map[string]interface{}

func (jsonField JSONB) Value() (driver.Value, error) {
	return json.Marshal(jsonField)
}

func (jsonField *JSONB) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &jsonField)
}

type User struct {
	ID       uint
	Username string
	Password string
}

var DB *gorm.DB

func InitializeDatabase() error {
	var err error

	dsn := "host=localhost user=postgres password=a dbname=web port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return err
}

func AddUser(username string, password string) error {
	tx := DB.Exec("SELECT * FROM users where username = ?;", username)
	if tx.RowsAffected != 0 {
		return errors.New(custom_error.USED_USERNAME)
	}

	DB.Exec("INSERT INTO users (username, password) VALUES(?, ?);",
		username, password)

	return nil
}

func GetUser(username string) (User, error) {
	var user User
	tx := DB.Raw("SELECT * FROM users where username = ?;", username).Scan(&user)
	if tx.RowsAffected != 1 {
		return User{}, errors.New(custom_error.INVALID_CREDENTIALS)
	}

	return user, nil
}

func GetBaskets(user_id uint) []Basket {
	var baskets []Basket
	DB.Raw("SELECT * FROM baskets WHERE user_id = ?;", user_id).Scan(&baskets)

	return baskets
}

func CreateBasket(user_id uint, data map[string]interface{}, state string) (uint, error) {
	if state != "COMPLETED" && state != "PENDING" {
		return 0, errors.New(custom_error.INVALID_STATE)
	}

	var id uint
	DB.Raw("INSERT INTO baskets (user_id, created_at, updated_at, data, state) VALUES(?, ?, ?, ?, ?) RETURNING id;",
		user_id, time.Now(), time.Now(), data, state).Scan(&id)

	return id, nil
}

func UpdateBasket(user_id uint, id uint, data map[string]interface{}, state string) error {
	if state != COMPLETED && state != PENDING {
		return errors.New(custom_error.INVALID_STATE)
	}

	var current_state string
	tx := DB.Raw("SELECT state FROM baskets WHERE id = ? AND user_id = ?;", id, user_id).Scan(&current_state)
	if tx.RowsAffected != 1 {
		return errors.New(custom_error.INVALID_ARGUMENTS)
	}
	if current_state == COMPLETED {
		return errors.New(custom_error.RESTRICTED_UPDATE)
	}

	tx = DB.Exec("UPDATE baskets SET updated_at = ?, data = ?, state = ? WHERE id = ? AND user_id = ?;",
		time.Now(), data, state, id, user_id)
	if tx.RowsAffected != 1 {
		return errors.New(custom_error.INVALID_ARGUMENTS)
	}

	return nil
}

func GetBasket(user_id uint, id uint) (Basket, error) {
	var basket Basket
	tx := DB.Raw("SELECT * FROM baskets WHERE id = ? AND user_id = ?;", id, user_id).Scan(&basket)
	if tx.RowsAffected != 1 {
		return Basket{}, errors.New(custom_error.INVALID_ARGUMENTS)
	}

	return basket, nil
}

func DeleteBasket(user_id uint, id uint) error {
	tx := DB.Exec("DELETE FROM baskets WHERE id = ? AND user_id = ?;", id, user_id)
	if tx.RowsAffected != 1 {
		return errors.New(custom_error.INVALID_ARGUMENTS)
	}

	return nil
}

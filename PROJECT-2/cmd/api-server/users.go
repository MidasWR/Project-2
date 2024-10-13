package api

import (
	"PROJECT-2/config"
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type Users struct {
	Id       int
	Username string
	Login    string
	Password string
}

func RegUser(config config.Config, username string, login string, password string) error {

	check, err := CheckLogin(login, config.DB)
	if err != nil {
		return err
	}
	if check == true {
		log.Println("User already logged in")
		return nil
	}
	check, err = CheckUsername(username, config.DB)
	if err != nil {
		return err
	}
	if check {
		log.Println("Name already used by users")
		return errors.New("name already used by users")
	}
	_, err = config.DB.Exec(`INSERT INTO users (username, login, password) VALUES ($1, $2, $3)`, username, login, password)
	if err != nil {
		return err
	}
	log.Println("User register access")
	return nil
}
func LoginUsers(db *sql.DB, login string, password string) string {
	var u Users
	u.Login = login
	u.Password = password
	var asd string
	if err := db.Ping(); err != nil {
		return "error3"
	}
	row := db.QueryRow(
		`SELECT  password FROM users WHERE login = $1`,
		u.Login,
	)
	err := row.Scan(&asd)
	if errors.Is(sql.ErrNoRows, err) {
		log.Println("User login not found")
		return "error4"
	}
	if asd != u.Password {
		log.Println("User login failed")
		return "error5"
	}
	token, err := CreateJWT()
	if err != nil {
		return "error6"
	}
	return token
}
func (u *Users) GetDataByLogin(login string, db *sql.DB) (*Users, error) {
	log.Println("1")
	err := db.QueryRow(
		`SELECT id,username,login,password FROM users WHERE login = $1`, login).
		Scan(
			&u.Id,
			&u.Username,
			&u.Login,
			&u.Password,
		)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CheckLogin(login string, db *sql.DB) (bool, error) {
	var asd int
	if err := db.Ping(); err != nil {
		return false, err
	}
	row := db.QueryRow(
		`SELECT  id FROM users WHERE login = $1`,
		login,
	)
	if row == nil {
		return false, nil
	}
	err := row.Scan(&asd)
	if errors.Is(sql.ErrNoRows, err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if asd > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
func CheckUsername(username string, db *sql.DB) (bool, error) {
	var asd int
	if err := db.Ping(); err != nil {
		return false, err
	}
	row := db.QueryRow(
		`SELECT  id FROM users WHERE login = $1`,
		username,
	)
	if row == nil {
		return false, nil
	}
	err := row.Scan(&asd)
	if errors.Is(sql.ErrNoRows, err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if asd > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
func CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Hour * 2160).Unix()
	AnsToken, err := token.SignedString(SignedKey)
	if err != nil {
		return "error", err
	}
	return AnsToken, nil
}

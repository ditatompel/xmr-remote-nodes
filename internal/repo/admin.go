package repo

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"xmr-remote-nodes/internal/database"

	"github.com/alexedwards/argon2id"
)

type Admin struct {
	Id           int    `db:"id"`
	Username     string `db:"username"`
	Password     string `db:"password"`
	LastactiveTs int64  `db:"lastactive_ts"`
	CreatedTs    int64  `db:"created_ts"`
}

type AdminRepo struct {
	db *database.DB
}

type AdminRepository interface {
	CreateAdmin(*Admin) (*Admin, error)
	Login(username string, password string) (*Admin, error)
}

func NewAdminRepo(db *database.DB) AdminRepository {
	return &AdminRepo{db}
}

func (repo *AdminRepo) CreateAdmin(admin *Admin) (*Admin, error) {
	if !validUsername(admin.Username) {
		return nil, errors.New("username is not valid, must be at least 4 characters long and contain only lowercase letters and numbers")
	}
	if !strongPassword(admin.Password) {
		return nil, errors.New("password is not strong enough, must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number and one special character")
	}
	hash, err := setPassword(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.Password = hash

	admin.CreatedTs = time.Now().Unix()

	if repo.isUsernameExists(admin.Username) {
		return nil, errors.New("username already exists")
	}

	query := `INSERT INTO tbl_admin (username, password, created_ts) VALUES (?, ?, ?)`
	_, err = repo.db.Exec(query, admin.Username, admin.Password, admin.CreatedTs)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (repo *AdminRepo) Login(username, password string) (*Admin, error) {
	query := `SELECT id, username, password FROM tbl_admin WHERE username = ? LIMIT 1`
	row, err := repo.db.Query(query, username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer row.Close()

	admin := Admin{}
	if row.Next() {
		err = row.Scan(&admin.Id, &admin.Username, &admin.Password)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid username or password")
	}

	match, err := checkPassword(admin.Password, password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if !match {
		return nil, errors.New("Invalid username or password")
	}

	update := `UPDATE tbl_admin SET lastactive_ts = ? WHERE id = ?`
	_, err = repo.db.Exec(update, time.Now().Unix(), admin.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &admin, nil
}

func (repo *AdminRepo) isUsernameExists(username string) bool {
	query := `SELECT id FROM tbl_admin WHERE username = ? LIMIT 1`
	row, err := repo.db.Query(query, username)
	if err != nil {
		return false
	}
	defer row.Close()
	return row.Next()
}

func setPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func checkPassword(hash, password string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func strongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	if !strings.ContainsAny(password, "0123456789") {
		return false
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}
	if !strings.ContainsAny(password, "!@#$%^&*()_+|~-=`{}[]:;<>?,./") {
		return false
	}
	return true
}

// No special character and unicode for username
func validUsername(username string) bool {
	if len(username) < 5 || len(username) > 20 {
		return false
	}

	// reject witespace, tabs, newlines, and other special characters
	if strings.ContainsAny(username, " \t\n") {
		return false
	}
	// reject unicode
	if strings.ContainsAny(username, "^\x00-\x7F") {
		return false
	}
	// reject special characters
	if strings.ContainsAny(username, "!@#$%^&*()_+|~-=`{}[]:;<>?,./ ") { // note last blank space
		return false
	}

	if !strings.ContainsAny(username, "abcdefghijklmnopqrstuvwxyz0123456789") {
		return false
	}
	return true
}

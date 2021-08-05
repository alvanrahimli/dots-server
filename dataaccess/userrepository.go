package dataaccess

import (
	"database/sql"
	"github.com/alvanrahimli/dots-server/models"
)

func FindUserById(db *sql.DB, userId int) (models.User, error) {
	user := models.User{}

	getUserRawQuery := `SELECT Id, Username, Email, Password FROM Users u WHERE u.Id = $1`
	row := db.QueryRow(getUserRawQuery, userId)
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func FindUserByEmail(db *sql.DB, email string) (models.User, error) {
	user := models.User{}

	getUserRawQuery := `SELECT Id, Username, Email, Password FROM Users u WHERE u.Email = $1`
	row := db.QueryRow(getUserRawQuery, email)
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}

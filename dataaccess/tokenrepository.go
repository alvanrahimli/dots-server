package dataaccess

import (
	"database/sql"
	"github.com/alvanrahimli/dots-server/models"
)

func AddToken(token *models.AuthToken, db *sql.DB) error {
	insertTokenQuery := `INSERT INTO AuthTokens (Token, ExpirationDate, UserId) VALUES (?, ?, ?)`
	statement, stmtErr := db.Prepare(insertTokenQuery)
	if stmtErr != nil {
		return stmtErr
	}

	_, dbErr := statement.Exec(token.Token, token.ExpirationDate, token.UserId)
	return dbErr
}

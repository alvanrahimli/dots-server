package utils

import (
	"crypto/sha256"
	"database/sql"
	"github.com/alvanrahimli/dots-server/models"
	"regexp"
	"strings"
	"time"
)

func ValidateName(name string) (string, bool) {
	isValid := true
	errors := make([]string, 0)

	if len(strings.TrimSpace(name)) < 4 {
		isValid = false
		errors = append(errors, "Username must be at least 4 characters long")
	}

	return strings.Join(errors, ", "), isValid
}

func ValidateEmail(email string) (string, bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	isValid := true
	errors := make([]string, 0)

	if !re.MatchString(email) {
		isValid = false
		errors = append(errors, "Invalid email")
	}

	return strings.Join(errors, ", "), isValid
}

func ValidatePassword(password string) (string, bool) {
	isValid := true
	errors := make([]string, 0)

	if len(strings.TrimSpace(password)) < 8 {
		isValid = false
		errors = append(errors, "Password must be at least 8 characters long")
	}

	return strings.Join(errors, ", "), isValid
}

func ValidatePackage(pack *models.Package) (string, bool) {
	isValid := true
	errors := make([]string, 0)

	if len(pack.Name) < 5 {
		isValid = false
		errors = append(errors, "Package name must be at least 5 characters long")
	}

	versionRe := regexp.MustCompile("^([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})")
	if !versionRe.MatchString(pack.Version) {
		isValid = false
		errors = append(errors, "Package version is invalid. Must be: x.x.x form")
	}

	return strings.Join(errors, ", "), isValid
}

func HashPassword(password string) string {
	result := sha256.Sum256([]byte(password))
	return string(result[:])
}

func ValidateToken(token string, db *sql.DB) (int, bool) {
	var userId int
	var expDateStr string
	// Get token
	getTokenQuery := `SELECT ExpirationDate, UserId FROM AuthTokens at WHERE at.Token = $1`
	row := db.QueryRow(getTokenQuery, token)
	if err := row.Scan(&expDateStr, &userId); err != nil {
		return 0, false
	}

	// Check expiration date
	expDate, timeErr := time.Parse(models.DatetimeLayout, expDateStr)
	if timeErr != nil {
		return 0, false
	}

	if time.Now().Before(expDate) {
		return userId, true
	}

	return 0, false
}

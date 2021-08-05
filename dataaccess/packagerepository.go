package dataaccess

import (
	"database/sql"
	"github.com/alvanrahimli/dots-server/models"
)

func AddPackage(pack *models.Package, db *sql.DB) (models.Package, error) {
	insertQuery := `INSERT INTO Packages (Name, Version, ArchiveName, UserId) VALUES (?, ?, ?, ?)`
	statement, stmtErr := db.Prepare(insertQuery)
	if stmtErr != nil {
		return models.Package{}, stmtErr
	}

	result, execErr := statement.Exec(pack.Name, pack.Version, pack.ArchiveName, pack.UserId)
	if execErr != nil {
		return models.Package{}, execErr
	}

	packId, lastErr := result.LastInsertId()
	if lastErr != nil {
		return models.Package{}, lastErr
	}

	addedPackage, getErr := GetPackage(int(packId), db)
	if getErr != nil {
		return models.Package{}, getErr
	}

	return addedPackage, nil
}

func GetPackage(packId int, db *sql.DB) (models.Package, error) {
	pack := models.Package{
		Id:      0,
		Name:    "",
		Version: "",
		UserId:  0,
	}

	getQuery := `SELECT Id, Name, Version, ArchiveName, UserId FROM Packages p WHERE p.Id = $1`
	row := db.QueryRow(getQuery, packId)
	if err := row.Scan(&pack.Id, &pack.Name, &pack.Version, &pack.ArchiveName, &pack.UserId); err != nil {
		return models.Package{}, err
	}

	return pack, nil
}

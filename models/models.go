package models

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}

type AuthToken struct {
	Id             int
	UserId         int
	Token          string
	ExpirationDate string
}

type Package struct {
	Id      int
	Name    string
	Version string

	UserId int
}

type App struct {
	Id      int
	Name    string
	Version string

	PackageId int
}

type RemoteAddr struct {
	Name string
	Url  string

	PackageId int
}

const (
	DatetimeLayout = "02/01/2006 15:04:05"
)

package models

import "os"

const (
	DatetimeLayout     = "02/01/2006 15:04:05"
	ArchivesFolderRoot = "./archives"
)

var RegistryDomain = os.Getenv("REGISTRY_DOMAIN")

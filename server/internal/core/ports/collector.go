package ports

import "database/sql"

type MattermostDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Stats() sql.DBStats
}

type MattermostAPI interface {
}

type MattermostPluginAPI interface {
}

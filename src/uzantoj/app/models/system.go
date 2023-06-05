package models


type System struct {
	Id                  int64    `db:"id"`
	LastUID             int      `db:"last_uid"`
	LastGID             int      `db:"last_gid"`
	Installed           bool     `db:"installed"`
	Version             int      `db:"version"`
}

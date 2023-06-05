package models

type DomainUser struct {
	Id              int64      `db:"id"`
	DomainId        int        `db:"domain_id"`
	UserId          int        `db:"user_id"`
	Created         string	   `db:"created,size:255"`
}

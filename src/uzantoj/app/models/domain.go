package models

import (
	"fmt"
	"github.com/revel/revel"
)

type Domain struct {
	Id              int64      `db:"id"`
	Gid             int        `db:"gid"`
	Name            string     `db:"name,size:255"`
	Created         string	   `db:"created,size:255"`
	Admins          []User     `db:"-"`
}


func (d *Domain) String() string {
	return fmt.Sprintf("%s", d.Name)
}

func (d *Domain) DomainNameForLDAP(ldap_root string) string {
	return "dc="+d.Name+","+ldap_root
}

func (d *Domain) UsersBranchForLDAP(ldap_root string) string {
	return "ou=users,"+d.DomainNameForLDAP(ldap_root)
}

func DomainCreate(name string, gid int) *Domain {
	domain := &Domain{
		Created: now(),
		Name: name,
		Gid: gid,
	}
	return domain
}

func (d *Domain) Validate(v *revel.Validation) {
	v.Domain(d.Name).Message("El nombre de dominio indicado no es válido.")
	v.MinSize(d.Name, 3).Message("El nombre es muy corto. Debe tener al menos 3 caracteres.")
	v.MaxSize(d.Name, 255).Message("El nombre es muy largo. Debe tener como máximo 255 caracteres.")	
}



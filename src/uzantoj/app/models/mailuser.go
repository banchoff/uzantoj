package models

import (
	"github.com/revel/revel"
	"regexp"
	"strings"
)

type MailUser struct {
	Uid             int
	Gid             int
	Username        string
	Name            string
	Lastname        string
	Domain          string
	Email           string
	UidName         string
}

func MailUserCreate(name, lastname, username, domain string) *MailUser {
	mailuser := &MailUser{
		Name: name,
		Lastname: lastname,
		Username: username,
		Email: username+"@"+domain,
		UidName: username+"@"+domain,
		Domain: domain,
	}
	return mailuser
}

func (m *MailUser) GetMyLDAPBranch(ldap_root string) string {
	ns := strings.Replace(m.Email, ".", "_", -1)
	ns = strings.Replace(ns, "@", "_", -1)	
	return "uid="+ns+",ou=users,dc="+m.Domain+","+ldap_root
}

func (m *MailUser) GetMyLDAPBranch2(ldap_root string) string {
	return "uid="+m.UidName+",ou=users,dc="+m.Domain+","+ldap_root
}


func (m *MailUser) Validate(v *revel.Validation) {
	v.Required(m.Name).Message("Debe indicar un nombre.")
	v.Match(m.Name, regexp.MustCompile("^[\\w\\s]*$")).Message("El nombre contiene caracteres inválidos.")
	v.MinSize(m.Name, 3).Message("El nombre es muy corto. Debe tener al menos 3 caracteres.")
	v.MaxSize(m.Name, 50).Message("El nombre es muy largo. Debe tener como máximo 50 caracteres.")

	v.Required(m.Lastname).Message("Debe indicar un apellido.")
	v.Match(m.Lastname, regexp.MustCompile("^[\\w\\s]*$")).Message("El apellido contiene caracteres inválidos.")
	v.MinSize(m.Lastname, 3).Message("El apellido es muy corto. Debe tener al menos 3 caracteres.")
	v.MaxSize(m.Lastname, 50).Message("El apellido es muy largo. Debe tener como máximo 50 caracteres.")
}

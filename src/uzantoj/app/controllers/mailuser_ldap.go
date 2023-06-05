package controllers

import (
	"github.com/revel/revel"
	"strconv"
	"github.com/go-ldap/ldap"
	"uzantoj/app/models"
	"errors"
)

func addMailUser(mailuser *models.MailUser, password string) bool {
	uid_   := getUID_(mailuser)
	ldap_root := revel.Config.StringDefault("ldap.root", "dc=nodomain")

	baseDN := mailuser.GetMyLDAPBranch(ldap_root)
	
	l := getLDAP()
	defer closeLDAP(l)

	// NOTA: el GID lo obtiene de la info del dominio guardada en la base de datos, no en el LDAP.
	//       El GID guardado en MySQL y en LDAP deben coincidir, de lo contrario se crearian usuarios para otro dominio.

	homeDir    := revel.Config.StringDefault("host.homedir", "/var/vdomains")+"/"+mailuser.Domain+"/users/"+uid_
	gidNumber  := strconv.Itoa(mailuser.Gid)
	uidNumber  := strconv.Itoa(mailuser.Uid)
	loginShell := "/bin/false"
	
	addReq := ldap.NewAddRequest(baseDN, []ldap.Control{})
	addReq.Attribute("objectClass",      []string{"posixAccount", "inetOrgPerson"})
	addReq.Attribute("uid",              []string{uid_, mailuser.Email})
	addReq.Attribute("givenName",        []string{mailuser.Name})
	addReq.Attribute("sn",               []string{mailuser.Lastname})
	addReq.Attribute("cn",               []string{mailuser.Username})
	addReq.Attribute("uidNumber",        []string{uidNumber})
	addReq.Attribute("gidNumber",        []string{gidNumber})
	addReq.Attribute("homeDirectory",    []string{homeDir})
	addReq.Attribute("mail",             []string{mailuser.Email})
	addReq.Attribute("loginShell",       []string{loginShell})

	if err := l.Add(addReq); err != nil {
		revel.AppLog.Error("LDAP: error adding user account:", addReq, err)
		return false
	}

	passwdModReq := ldap.NewPasswordModifyRequest(baseDN, "", password)
	if _, err := l.PasswordModify(passwdModReq); err != nil {
		revel.AppLog.Error("LDAP :: Error: failed to modify password: %v", err)
		// TODO: Podemos intentar borrar la rama creada
		return false
	} 
	
	if !addMailuserDirectory(mailuser) {
		revel.AppLog.Error("Error creando el directorio para el usuario.")
		// TODO: Podemos intentar borrar la rama creada
		return false
	} 

	return true
}

// TODO: Usar Ansible para hacer un backup y dejarlo en el servidor de backups.
// Esto se debe hacer con un CRON (soportado por Revel).
func deleteMailUser(uid int, domain string) bool {
	mailuser, err := getMailUser(uid, domain)
	if err != nil {
		revel.AppLog.Error("LDAP: error en deleteMailUser: no se encuentra el usuario", err)
		return false
	}
	
	ldap_root := revel.Config.StringDefault("ldap.root", "dc=nodomain")
	baseDN := mailuser.GetMyLDAPBranch2(ldap_root)

	l := getLDAP()
	defer closeLDAP(l)
	
	delReq := ldap.NewDelRequest(baseDN, []ldap.Control{})

	if err = l.Del(delReq); err != nil {
		revel.AppLog.Error("LDAP: error en deleteMailUser: no se puede borrar el usuario", err)
		return false
	}

	if !deleteMailuserDirectory(&mailuser) {
		revel.AppLog.Error("Error en deleteMailUser: no se puede borrar el directorio del usuario, pero se borro el usuario del LDAP..")
		return false
	}
	
	return true
}

func editMailUser(mailuser *models.MailUser) bool {
	var err error
	ldap_root := revel.Config.StringDefault("ldap.root", "dc=nodomain")
	baseDN := mailuser.GetMyLDAPBranch2(ldap_root)

	l := getLDAP()
	defer closeLDAP(l)
	
	modReq := ldap.NewModifyRequest(baseDN, []ldap.Control{})

	modReq.Replace("givenName", []string{mailuser.Name})
	modReq.Replace("sn",        []string{mailuser.Lastname})
	
	if err = l.Modify(modReq); err != nil {
		revel.AppLog.Error("LDAP: error modifying user account:", modReq, err)
	}
	
	return (err == nil)
}

func changePasswordMailUser(uid int, password, domain string) bool {

	var err error
	mailuser, _ := getMailUser(uid, domain)
	ldap_root := revel.Config.StringDefault("ldap.root", "dc=nodomain")
	baseDN := mailuser.GetMyLDAPBranch2(ldap_root)

	l := getLDAP()
	defer closeLDAP(l)

	passwdModReq := ldap.NewPasswordModifyRequest(baseDN, "", password)
	if _, err = l.PasswordModify(passwdModReq); err != nil {
		revel.AppLog.Error("LDAP :: Error: failed to modify password: %v", err)
	}

	return (err == nil)
}

// TODO: Implementar la paginacion de los resultados.
func listMailUsers(domain_name string) []models.MailUser {
	var mailusers []models.MailUser

	l := getLDAP()
	defer closeLDAP(l)

	ldap_root := revel.Config.StringDefault("ldap.root", "dc=nodomain")
	domain, err := getDomainByName(domain_name)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}
	
	baseDN := domain.UsersBranchForLDAP(ldap_root)
	
	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeSingleLevel, 0, 0, 0, false, "(objectClass=*)", []string{}, []ldap.Control{})
	
	result, err := l.Search(searchReq)
	if err != nil {
		revel.AppLog.Error("LDAP :: Error: failed to query LDAP: %w", err)
	}

	if result != nil {
		for _, entry := range result.Entries {
			uid, _ := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
			gid, _ := strconv.Atoi(entry.GetAttributeValue("gidNumber"))
			mailusers = append(mailusers, models.MailUser{
				Uid: uid,
				Gid: gid,
				Username: entry.GetAttributeValue("cn"),
				Name: entry.GetAttributeValue("givenName"),
				Lastname: entry.GetAttributeValue("sn"),
				Domain: domain.Name,
				Email: entry.GetAttributeValue("mail"),
				UidName: entry.GetAttributeValue("uid"),
			})		
		}
	}
	return mailusers
}

func getMailUser(uid int, domain_name string) (models.MailUser, error) {
	var err error
	mailuser :=  models.MailUser{}
	found := false
	err = nil
	mailusers :=  listMailUsers(domain_name)
	for i := 0; i < len(mailusers); i++ {
		if mailusers[i].Uid == uid {
			mailuser = mailusers[i]
			found = true
			break
		}
	}
	if !found {
		err = errors.New("No se encontró el usuario.")
	}
	return mailuser, err
}


func hasMailUsers(domain string) bool {
	mailusers := listMailUsers(domain)
	return len(mailusers) > 0
}

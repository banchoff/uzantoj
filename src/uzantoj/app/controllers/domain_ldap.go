package controllers

import (
	"github.com/revel/revel"
	"uzantoj/app/models"
	"strconv"
	"github.com/go-ldap/ldap"
	"strings"
)

func addDomainToLDAP(domain *models.Domain) (*models.Domain, error) {
	gid := strconv.Itoa(domain.Gid)
	uid := strconv.Itoa(getNextUID())
	ldap_root := revel.Config.StringDefault("ldap.root", "dc=nodomain")
	homeDir  := revel.Config.StringDefault("host.homedir", "/var/vdomains/")

	baseDN_1 := domain.DomainNameForLDAP(ldap_root)
	baseDN_2 := domain.UsersBranchForLDAP(ldap_root)

	l := getLDAP()
	defer closeLDAP(l)

	addReq_1 := ldap.NewAddRequest(baseDN_1, []ldap.Control{})
	addReq_2 := ldap.NewAddRequest(baseDN_2, []ldap.Control{})
	// Rama del dominio
	addReq_1.Attribute("objectClass",   []string{"posixAccount", "organization", "top", "dcObject"})
	addReq_1.Attribute("cn",            []string{domain.Name})
	addReq_1.Attribute("gidNumber",     []string{gid})
	addReq_1.Attribute("homeDirectory", []string{homeDir})
	addReq_1.Attribute("uidNumber",     []string{uid})
	addReq_1.Attribute("uid",           []string{domain.Name})
	addReq_1.Attribute("o",             []string{"nodomain"})	
	// Rama del dominio + ou=users
	addReq_2.Attribute("objectClass",   []string{"organizationalUnit", "top"})
	addReq_2.Attribute("ou",            []string{"users"})	
	
	if err := l.Add(addReq_1); err != nil {
		revel.AppLog.Error("LDAP: error adding domain:", addReq_1, err)
		return domain, err
	}

	if err := l.Add(addReq_2); err != nil {
		// TODO: Podemos intentar borrar la rama creada.
		revel.AppLog.Error("LDAP: error adding ou=users,$domain:", addReq_2, err)
		return domain, err
	}
	return domain, nil
}

func deleteDomain(domain_name string) bool {
	ldap_root := revel.Config.StringDefault("ldap.root", "dc=nodomain")
	domain, err := getDomainByName(domain_name)
	if err != nil {
		revel.AppLog.Error("Error -", err)
		return false
	}
	if hasMailUsers(domain.Name) {
		revel.AppLog.Error("Hay usuarios de email en el dominio. Para borrar el dominio no debe haber usuarios.")
		return false
	}

	// Asumimos que el dominio existe.	
	baseDN_1 := domain.UsersBranchForLDAP(ldap_root)
	baseDN_2 := domain.DomainNameForLDAP(ldap_root)
	
	l := getLDAP()
	defer closeLDAP(l)
	
	delReq_1 := ldap.NewDelRequest(baseDN_1, []ldap.Control{})
	delReq_2 := ldap.NewDelRequest(baseDN_2, []ldap.Control{})

	if err := l.Del(delReq_1); err != nil {
		revel.AppLog.Error("LDAP: error en deleteDomain: no se puede borrar ou=users,$domain", err)
		return false
	}

	if err := l.Del(delReq_2); err != nil {
		revel.AppLog.Error("LDAP: error en deleteDomain: no se puede borrar el dominio", err)
		return false
	}

	if !deleteDomainDirectory(domain.Name) {
		revel.AppLog.Error("No se pudo borrar el directorio del dominio, pero se borro el dominio del LDAP. Recomendacion: borrar el directorio a mano.")
		return false
	}
	
	return true
} 

func getDomainNamesLDAP() []string {
	var domains []string
	l := getLDAP()
	defer closeLDAP(l)

	baseDN := revel.Config.StringDefault("ldap.root", "dc=nodomain")
	
	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeSingleLevel, 0, 0, 0, false, "(objectClass=*)", []string{}, []ldap.Control{})	
	result, err := l.Search(searchReq)
	if err != nil {
		revel.AppLog.Error("LDAP :: Error: failed to query LDAP: %w", err)
	}
	
	for _, entry := range result.Entries {
		tmp := strings.Split(entry.DN, "=")
		// Forma: tmp[1] === testing.com.ar,ou
		tmp2 := strings.Split(tmp[1], ",")
		// Para saber si es un dominio:
		if strings.Contains(tmp2[0], ".") {
			domains = append(domains, tmp2[0])
		}
	}
	return domains
}

func validateDir(directory string) bool {
	if strings.Contains(directory, "..") {
		return false
	}
	// TODO: mas tests....
	return true
}

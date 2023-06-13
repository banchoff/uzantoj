package controllers

import (
	"github.com/revel/revel"
	"uzantoj/app/models"
)

func saveDomain(aDomain *models.Domain) (*models.Domain, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	var err error
	
	aDomain.Created = now()
	if (aDomain.Id > 0) {
		_, err = DBMAP.Update(aDomain)
	} else {
		err = DBMAP.Insert(aDomain)
	}
	if err != nil {
		revel.AppLog.Error("Domain.saveDomain - Error - ", err)
	}
	return aDomain, err
}

func domainAssignedTo(domain_id, user_id int) (models.Domain, error) {
	var domain models.Domain
	DBMAP := getDB()
	defer closeDB(DBMAP)
	err := DBMAP.SelectOne(&domain, `SELECT domains.* 
                                         FROM domains, domains_users 
                                         WHERE domains.id = domains_users.domain_id 
                                         AND domains_users.user_id = ? and domains.id = ?`, user_id, domain_id)
	return domain, err
}

func domainsAssignedTo(user_id int) ([]models.Domain, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	var domains []models.Domain
	_, err := DBMAP.Select(&domains, `SELECT domains.* 
                                          FROM domains, domains_users  
                                          WHERE domains.id = domains_users.domain_id
                                          AND domains_users.user_id=?
                                          ORDER BY id`, user_id)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}
	return domains, err
}


func hasAdmins(domain_id int) bool {
	users, _ := getAdminsForDomain(domain_id)
	return (len(users) > 0)
}

func getNextGID() int {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	system := models.System{}

	err := DBMAP.SelectOne(&system, `SELECT * 
                                         FROM system 
                                         WHERE id = 1`)
	if err != nil {
		revel.AppLog.Error("Error getNextGID: ", err)
		return -1
	}
	gid := system.LastGID
	system.LastGID = system.LastGID + 1
	_, err = DBMAP.Update(&system)
	if err != nil {
		revel.AppLog.Error("Error getNextGID: ", err)
		return -1
	}
	return gid
}

func getDomainById(domain_id int) (models.Domain, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	domain := models.Domain{}
	err := DBMAP.SelectOne(&domain, `SELECT * 
                                         FROM domains 
                                         WHERE id=?`, domain_id)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
		domain.Id = -1
	}
	return domain, err
}

func getDomainByName(domain_name string) (models.Domain, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	domain := models.Domain{}
	err := DBMAP.SelectOne(&domain, `SELECT * 
                                         FROM domains 
                                         WHERE name=?`, domain_name)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
		domain.Id = -1
	}
	return domain, err
}

func domainExistsByName(domain_name string) bool {
	_, err := getDomainByName(domain_name)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}
	return (err == nil)
}

func assignUserToDomain(aUser models.User, aDomain models.Domain) (models.DomainUser, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	du := &models.DomainUser{
		DomainId: int(aDomain.Id),
		UserId: int(aUser.Id),
		Created: now(),
	}
	err := DBMAP.Insert(du)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}
	return *du, err
}

func unAssignUserToDomain(user_id, domain_id int) error {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	_, err := DBMAP.Exec(`DELETE FROM domains_users 
                              WHERE domain_id=? 
                              AND user_id=?`, domain_id, user_id)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}
	return err
}

func deleteDomainById(domain_id int) error {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	_, err := DBMAP.Exec(`DELETE FROM domains 
                              WHERE id=?`, domain_id)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}
	return err
}

func getAllDomains() ([]models.Domain, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	var domains []models.Domain
	_, err := DBMAP.Select(&domains, `SELECT * 
                                          FROM domains 
                                          ORDER BY id`)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}
	return domains, err
}


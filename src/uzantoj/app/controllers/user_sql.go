package controllers

import (
	"github.com/revel/revel"
	"uzantoj/app/models"
)

func hasDomains(user_id int) bool {
	domains, err := domainsAssignedTo(user_id)
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}
	return (len(domains) > 0)
}

func getUserById(user_id int) (models.User, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	user := models.User{}
	err := DBMAP.SelectOne(&user, `SELECT * 
                                       FROM users 
                                       WHERE id=?`, user_id)
	if err != nil {
		revel.AppLog.Error("User.getUserById - Error - ", err)
		user.Id = -1
	}
	return user, err
}

func getUserByUsername(username string) (models.User, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	user := models.User{}
	err := DBMAP.SelectOne(&user, `SELECT * 
                                       FROM users 
                                       WHERE username=?`, username)
	if err != nil {
		revel.AppLog.Error("User.getUserByUsername - Error - ", err)
		user.Id = -1
	}
	return user, err
}

func getAllUsers() ([]models.User, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	var users []models.User
	_, err := DBMAP.Select(&users, `SELECT * 
                                        FROM users 
                                        ORDER BY id`)
	if err != nil {
		revel.AppLog.Error("User.getAllUsers - Error - ", err)
	}
	return users, err
}

func getAllUsersByRole(aRole string) ([]models.User, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	var users []models.User
	_, err := DBMAP.Select(&users, `SELECT * 
                                        FROM users
                                        WHERE role = ? 
                                        ORDER BY id`, aRole)
	if err != nil {
		revel.AppLog.Error("User.getAllUsersByRole - Error - ", err)
	}
	return users, err
}

func userExists(user_id int) bool {
	_, err := getUserById(user_id)
	if err != nil {
		revel.AppLog.Error("Lib.UserExists - Error - ", err)
	}
	return err == nil
}

func deleteUserById(user_id int) error {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	_, err := DBMAP.Exec(`DELETE FROM users 
                              WHERE id=?`, user_id)
	if err != nil {
		revel.AppLog.Error("User.Delete - Error borrando el usuario - ", err)
	}
	return err
}

func getAdminsForDomain(domain_id int) ([]models.User, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	var users []models.User
	_, err := DBMAP.Select(&users, `SELECT users.* 
                                       FROM users, domains_users 
                                       WHERE users.id = domains_users.user_id 
                                       AND domain_id=?`, domain_id)
	return users, err
}

func saveUser(aUser *models.User) (*models.User, error) {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	var err error
	
	aUser.Created = now()
	if (aUser.Id > 0) {
		_, err = DBMAP.Update(aUser)
	} else {
		err = DBMAP.Insert(aUser)
	}
	if err != nil {
		revel.AppLog.Error("User.saveUser - Error - ", err)
	}
	return aUser, err
}

package controllers

import (
	"github.com/revel/revel"
	"uzantoj/app/models"
	"github.com/go-gorp/gorp"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
	"github.com/go-ldap/ldap"
	"strings"
	"os/exec"
)

func now() string {
	current_time := time.Now()
	tmp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		current_time.Year(), current_time.Month(), current_time.Day(),
		current_time.Hour(), current_time.Minute(), current_time.Second())
	return tmp
}

func closeDB(d *gorp.DbMap) {
	d.Db.Close()
}

func closeLDAP(l *ldap.Conn) {
	l.Close()
}

func getDB() *gorp.DbMap {
	db_driver      := revel.Config.StringDefault("db.driver",   "mysql")
	db_username    := revel.Config.StringDefault("db.username", "username")
	db_password    := revel.Config.StringDefault("db.password", "password")
	db_server      := revel.Config.StringDefault("db.server",   "127.0.0.1")
	db_port        := revel.Config.StringDefault("db.port",     "3306")
	db_name        := revel.Config.StringDefault("db.name",     "localdb")
	
	connect_string := db_username+":"+db_password+"@tcp("+db_server+":"+db_port+")/"+db_name
	db, err := sql.Open(db_driver, connect_string)
	if err != nil {
		revel.AppLog.Error("DB Error", err)
	}
	
	DBMAP := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	DBMAP.AddTableWithName(models.User{},               "users").SetKeys(true, "Id") // Id es el campo del struct, no de la table.
	DBMAP.AddTableWithName(models.System{},            "system").SetKeys(true, "Id")
	DBMAP.AddTableWithName(models.Domain{},           "domains").SetKeys(true, "Id")
	DBMAP.AddTableWithName(models.DomainUser{}, "domains_users").SetKeys(true, "Id")
	return DBMAP
}

func getLDAP() *ldap.Conn {
	ldap_server   := revel.Config.StringDefault("ldap.server",   "localhost")
	ldap_port     := revel.Config.StringDefault("ldap.port",     "389")
	ldap_bind     := revel.Config.StringDefault("ldap.bind",     "root")
	ldap_password := revel.Config.StringDefault("ldap.password", "password")
	
	l, err := ldap.DialURL("ldap://"+ldap_server+":"+ldap_port)
	if err != nil {
		revel.AppLog.Error("LDAP :: Dial error", err)
	}

	err = l.Bind(ldap_bind, ldap_password)
	if err != nil {
		revel.AppLog.Error("LDAP :: Bind error", err)
	}
	
	return l
}

func getUID_(mailuser *models.MailUser) string {
	ns := strings.Replace(mailuser.Email, ".", "_", -1)
	ns = strings.Replace(ns, "@", "_", -1)
	return ns
}

func getNextUID() int {
	DBMAP := getDB()
	defer closeDB(DBMAP)
	system := models.System{}

	err := DBMAP.SelectOne(&system, `SELECT * 
                                         FROM system 
                                         WHERE id = 1`)
	if err != nil {
		revel.AppLog.Error("Error getNextUID: ", err)
		return -1
	}

	uid := system.LastUID
	system.LastUID = system.LastUID + 1
	_, err = DBMAP.Update(&system)
	if err != nil {
		revel.AppLog.Error("Error getNextUID: ", err)
		return -1
	}

	return uid
}

func runCommand(commandName, dirName string) bool {
	hostName := revel.Config.StringDefault("host.name", "multidomain.local.test")
	hostUser := revel.Config.StringDefault("host.user", "uzantoj_create_domain.test")

	if !validateDir(dirName) {
		return false
	}
	
	command := exec.Command(commandName, dirName, hostName, hostUser)
	stdout, err := command.Output()

	if err != nil {
		revel.AppLog.Error("Error ejecutando: ", err)
		revel.AppLog.Error("Salida del comando con error: "+string(stdout))
	}
	return (err == nil)
}

// TODO: lib.sendMetric
// Envia una metrica a Graphite, InfluxDB o similar
// La metrica tiene la forma uzantoj.INSTALL-NAME.ALGO
func sendMetric(metric, value string) bool {
	return true
}

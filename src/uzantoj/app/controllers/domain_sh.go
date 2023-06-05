package controllers

import (
	"uzantoj/app/models"
	"github.com/revel/revel"
)

func executeCommand(cmd, rPath string) bool {
	baseDir := revel.Config.StringDefault("host.homedir", "/var/vdomains")
	dirName := baseDir+"/"+rPath
	return runCommand(cmd, dirName)
}

func addDomainDirectory(domain string) bool {
	// baseDir := revel.Config.StringDefault("host.homedir", "/var/vdomains")
	// dirName := baseDir+"/"+domain
	// return runCommand("bin/create_domain.sh", dirName)
	return executeCommand("bin/create_domain.sh", domain)
}

func addMailuserDirectory(mailuser *models.MailUser) bool {
	// baseDir := revel.Config.StringDefault("host.homedir", "/var/vdomains")
	// dirName := baseDir+"/"+mailuser.Domain+"/users/"+getUID_(mailuser)
	// return runCommand("bin/create_maildir.sh", dirName)
	return executeCommand("bin/create_maildir.sh", mailuser.Domain+"/users/"+getUID_(mailuser))
}

func deleteDomainDirectory(domain string) bool {
	// baseDir := revel.Config.StringDefault("host.homedir", "/var/vdomains")
	// dirName := baseDir+"/"+domain
	// return runCommand("bin/delete_domain.sh", dirName)	
	return executeCommand("bin/delete_domain.sh", domain)
}

func deleteMailuserDirectory(mailuser *models.MailUser) bool {
	// baseDir := revel.Config.StringDefault("host.homedir", "/var/vdomains")
	// dirName := baseDir+"/"+mailuser.Domain+"/users/"+getUID_(mailuser)
	// return runCommand("bin/delete_maildir.sh", dirName)
	return executeCommand("bin/delete_maildir.sh", mailuser.Domain+"/users/"+getUID_(mailuser))
}


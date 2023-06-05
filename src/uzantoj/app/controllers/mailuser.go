package controllers

import (
	"github.com/revel/revel"
	"uzantoj/app/models"
)

type MailUser struct {
	App
}

func (c MailUser) getMyDomainByDomainId(domain_id int) (models.Domain, error) {
	var domain models.Domain
	var err error
	
	if c.getMyRole() == "ADMIN" {
		domain, err = getDomainById(domain_id)
	} else {
		domain, err = domainAssignedTo(domain_id, c.getMyUID())
	}
	if err != nil {
		revel.AppLog.Error("Error - ", err)
	}

	return domain, err
}

func (c MailUser) validatePassword(password string) {
	c.Validation.MinSize(password, 10).Message("La contraseña es muy corta. Debe tener al menos 10 caracteres.")
	c.Validation.Required(password).Message("Debe indicar una contraseña.")
}

func (c MailUser) Index(id int) revel.Result {
	var mailusers []models.MailUser
	domain, err := c.getMyDomainByDomainId(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el dominio.")
		return c.RenderTemplate("App/Error.html")
	}	
	mailusers = listMailUsers(domain.Name)
	return c.Render(mailusers, domain)
}

func (c MailUser) Add(id int, username, name, lastname, password string) revel.Result {

	if c.getMyRole() != "ADMIN" && !isMyDomain(id, c.getMyUID()) {
		c.Flash.Error("El usuario no tiene permiso.")
		return c.RenderTemplate("App/Error.html")
	}

	domain, err := c.getMyDomainByDomainId(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el dominio.")
		return c.RenderTemplate("App/Error.html")
	}

      	if (c.Request.Method == "GET") {
		return c.Render(domain)
	}

	mailuser := models.MailUserCreate(name, lastname, username, domain.Name)
	mailuser.Validate(c.Validation)
	c.validatePassword(password)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/mailuser/add/%d/", id)
	}
	
	mailuser.Uid =  getNextUID()
	mailuser.Gid = domain.Gid

	if !addMailUser(mailuser, password) {
		c.Flash.Error("Error agregando el usuario en el LDAP")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/mailuser/add/%d/", id)
	}
	c.Flash.Success("Usuario agregado")
	return c.Redirect("/mailuser/view/%d/%d", domain.Id, mailuser.Uid)	
}

func (c MailUser) Delete(id, uid int) revel.Result {
	domain, err := c.getMyDomainByDomainId(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el dominio.")
		return c.RenderTemplate("App/Error.html")
	}

	if !deleteMailUser(int(uid), domain.Name) {
		c.Flash.Error("Error borrando el usuario del LDAP.")
		return c.RenderTemplate("App/Error.html")
	}
	return c.Render()
}

func (c MailUser) View(id, uid int) revel.Result {
	domain, err := c.getMyDomainByDomainId(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el dominio.")
		return c.RenderTemplate("App/Error.html")
	}
	
	mailuser, err := getMailUser(uid, domain.Name)
	if err != nil {
		c.Flash.Error("Error obteniendo el usuario del LDAP.")
		return c.RenderTemplate("App/Error.html")
	}
	return c.Render(mailuser, domain)
}

func (c MailUser) ChangePassword(id, uid int, password, password2 string) revel.Result {
      	if (c.Request.Method == "GET") {
		return c.Render(id, uid)
	}
	
	domain, err := c.getMyDomainByDomainId(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el dominio.")
		return c.RenderTemplate("App/Error.html")
	}
	c.validatePassword(password)
	c.Validation.Required(password == password2).Message("Las contraseñas deben ser iguales.")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/mailuser/change_password/%d/%d", id, uid)
	}

	if !changePasswordMailUser(uid, password, domain.Name) {
		c.Flash.Error("Error cambiando la contraseña del usuario en el LDAP.")
		//revel.AppLog.Error("MailUser.ChangePassword - Error cambiando la contraseña del usuario en el LDAP - ", err)
		return c.RenderTemplate("App/Error.html")
	}
	c.Flash.Success("Se modificó la contraseña")
	return c.Redirect("/mailuser/view/%d/%d", id, uid)
}

func (c MailUser) Edit(id, uid int, name, lastname string) revel.Result {

	domain, err := c.getMyDomainByDomainId(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el dominio.")
		return c.RenderTemplate("App/Error.html")
	}

	mailuser, err := getMailUser(uid, domain.Name)
	if err != nil {
		c.Flash.Error("Error obteniendo el usuario del LDAP.")
		//revel.AppLog.Error("MailUser.Edit - Error cargando el usuario- ", err)
		return c.RenderTemplate("App/Error.html")
	}
	
      	if (c.Request.Method == "GET") {
		return c.Render(domain, mailuser)
	}

	mailuser.Validate(c.Validation)	
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/mailuser/edit/%d/%d", id, uid)
	}

	mailuser.Name = name
	mailuser.Lastname = lastname
	if !editMailUser(&mailuser) {
		c.Flash.Error("Error actualizando el usuario en el LDAP.")
		//revel.AppLog.Error("MailUser.Edit - Error actualizando el usuario en el LDAP - ", err)
		return c.RenderTemplate("App/Error.html")
	}
	return c.Redirect("/mailuser/view/%d/%d", id, uid)
}


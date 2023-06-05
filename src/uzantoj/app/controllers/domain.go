package controllers

import (
	"github.com/revel/revel"
	"uzantoj/app/models"
)

type Domain struct {
	App
}

func (c Domain) Index() revel.Result {
	var domains         []models.Domain
	var domains_in_ldap []string
	var err               error

	// Logueamos si hay dominios en LDAP que no hay en la BD.
	domains_in_ldap = getDomainNamesLDAP()	
	for i := 0; i < len(domains_in_ldap); i++ {
		if !domainExistsByName(domains_in_ldap[i]) {
			revel.AppLog.Warn("Domain.Index: se encontró un dominio que no existe en la base de datos: ", domains_in_ldap[i])
		}
	}

	if c.getMyRole() == "ADMIN" {
		domains, err = getAllDomains()
	} else {
		domains, err = domainsAssignedTo(c.getMyUID())
	}
	if err != nil {
		c.Flash.Error("Error cargando los dominios.")
		// revel.AppLog.Error("Domain.Index - Error cargando dominio - ", err)
	}
	
	for i := 0; i < len(domains); i++ {
		domains[i].Admins, err = getAdminsForDomain(int(domains[i].Id))
		if err != nil {
			c.Flash.Error("Error cargando los admins para el dominio.", domains[i].Name)
			//revel.AppLog.Error("Domain.View - Error cargando relacion dominio X usuario - ", err)
		}
	}
	return c.Render(domains)
}

func (c Domain) AddAdmin(id, user_id int) revel.Result {
	
	users, err := getAllUsersByRole("USER")
	if err != nil {
		c.Flash.Error("Error cargando la lista de usuarios que pueden administrar dominios")
		//revel.AppLog.Error("Error - ", err)
		return c.RenderTemplate("App/Error.html")
	}
	
	domain, err := getDomainById(id)
	if err != nil {
		c.Flash.Error("Error cargando el dominio.")
		//revel.AppLog.Error("Domain.View - Error - ", err)
		return c.RenderTemplate("App/Error.html")
	}
	
      	if (c.Request.Method == "GET") {	
		return c.Render(users, domain)
	}

	user, err := getUserById(user_id)
	_, err = assignUserToDomain(user, domain)
	if err != nil {
		c.Flash.Error("Error asignando el dominio al usuario.")
		c.Validation.Keep()
		return c.Render(users, domain)
	}
	c.Flash.Success("Dominio agregado")
	return c.Redirect("/domain/view/%d", id)
}

func (c Domain) DeleteAdmin(id, uid int) revel.Result {

	err := unAssignUserToDomain(uid, id)
	if err != nil {
		c.Flash.Error("Error des-asignando el dominio al usuario.")
		revel.AppLog.Error("Error - ", err)
		return c.RenderTemplate("App/Error.html")
	}
	return c.Render()
}

func (c Domain) Add(name string) revel.Result {
      	if (c.Request.Method == "GET") {
		return c.Render()
	}

	domain := models.DomainCreate(name, getNextGID())
	domain.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Domain.Add)
	}

	_, err := addDomainToLDAP(domain)
	if err != nil {
		c.Flash.Error("Error agregando el dominio al LDAP.")
		return c.RenderTemplate("App/Error.html")
	}

	if !addDomainDirectory(domain.Name) {
		c.Flash.Error("Error creando el directorio para los mailboxes.")
		return c.RenderTemplate("App/Error.html")
	}

	_, err = saveDomain(domain)
	if err != nil {
		c.Flash.Error("Error agregando el dominio a la base de datos")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Domain.Add)
	}
	c.Flash.Success("Dominio agregado")
	return c.Redirect("/domain/view/%d", domain.Id)	
}

func (c Domain) Delete(id int) revel.Result {
	domain, err := getDomainById(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el dominio.")
		return c.RenderTemplate("App/Error.html")
	}

	if hasMailUsers(domain.Name) {
		c.Flash.Error("No se puede borrar un dominio que tiene usuarios de email. Borre primero los usuarios.")
		return c.RenderTemplate("App/Error.html")
	}
	if hasAdmins(id) {
		c.Flash.Error("No se puede borrar un dominio que está asignado a un administrador.")
		return c.RenderTemplate("App/Error.html")
	}
	
	if !deleteDomain(domain.Name) {
		c.Flash.Error("Error al borrar el dominio del LDAP.")
		return c.RenderTemplate("App/Error.html")
	}

	err = deleteDomainById(id)
	if err != nil {
		c.Flash.Error("Error al borrar el dominio de la BD.")
		return c.RenderTemplate("App/Error.html")
	}

	return c.Render()
}


func (c Domain) Users() revel.Result {
	return c.Render()
}


func (c Domain) View(id int) revel.Result {

	domain  :=  models.Domain{}
	var err           error
		
	if c.getMyRole() == "ADMIN" {
		domain, err = getDomainById(id)
	} else {
		domain, err = domainAssignedTo(id, c.getMyUID())
	}
	if err != nil {
		c.Flash.Error("Error obteniendo el dominio.")
		return c.RenderTemplate("App/Error.html")
	}

	domain.Admins, err = getAdminsForDomain(int(domain.Id))
	if err != nil {
		c.Flash.Error("Error obteniendo los administradores para el dominio.")
		revel.AppLog.Error("Domain.View - Error cargando relacion dominio X usuario - ", err)
		return c.RenderTemplate("App/Error.html")
	}

	return c.Render(domain)
}

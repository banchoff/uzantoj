package controllers

import (
	"github.com/revel/revel"
	"uzantoj/app/models"
)

type User struct {
	App
}

func (c User) Index() revel.Result {
	return c.Render()
}

func (c User) Login(username, password string) revel.Result {
      	if (c.Request.Method == "GET") {
		return c.Render()
	}

	user, _ := getUserByUsername(username)
	user.ValidatePasswords(c.Validation, password)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(User.Login)
	}

	c.Session["username"] = user.Username
	c.Session["role"] = user.Role
	c.Session["uid"] = user.Id
	return c.Redirect("/")
}

func (c User) Profile() revel.Result {
	myId2, _ := c.Session.Get("uid")
	myId := int(myId2.(float64))
	return c.Redirect("/user/view/%d", myId)
	
}

func (c User) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.RenderTemplate("User/LoggedOut.html")
}

func (c User) SystemUsers() revel.Result {
	users, err := getAllUsers()
	if err != nil {
		c.Flash.Error("Error cargando los usuarios del sistema.")
	}
	return c.Render(users)
}

func (c User) ChangePassword(id int, password, password2 string) revel.Result {
	if (c.getMyRole() != "ADMIN") && (id != c.getMyUID()) {
		c.Flash.Error("No es posible cambiar la contraseña de otro usuario.")
		return c.RenderTemplate("App/InvalidUser.html")
	}

      	if (c.Request.Method == "GET") {
		return c.Render()
	}
	
	user, err := getUserById(id)
	user.Validate(c.Validation)
	c.Validation.Required(password == password2).Message("Las contraseñas no coinciden.")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/user/change_password/%d", id)
	}

	user.ChangePassword(password)
	_, err = saveUser(&user)
	if err != nil {
		c.Flash.Error("Hubo un error al guardar las modificaciones del usuario en la base de datos.")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/user/change_password/%d", id)
	}
	c.Flash.Success("Se modificó la contraseña")
	return c.Redirect("/user/view/%d", user.Id)
} 

func (c User) Add(name, lastname, username, password, password2, email, role string) revel.Result {
	if (c.getMyRole() != "ADMIN") {
		c.Flash.Error("No tiene permisos para realizar esta operación.")
		return c.RenderTemplate("App/InvalidUser.html")
	}

      	if (c.Request.Method == "GET") {
		return c.Render()
	}

	user := models.UserCreate(name, lastname, username, password, role, email)
	c.Validation.Required(password == password2).Message("Las contraseñas no coinciden.")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(User.Add)
	}

	_, err := saveUser(user)
	if err != nil {
		c.Flash.Error("Error agregando el usuario a la base de datos")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(User.Add)
	}
	c.Flash.Success("Usuario agregado")
	return c.Redirect("/user/view/%d", user.Id)	
}

func (c User) Edit(id int, name, lastname, email, role string) revel.Result {
	if (c.getMyRole() != "ADMIN") && (id != c.getMyUID()) {
		c.Flash.Error("No puede editar otros usuarios.")
		return c.RenderTemplate("App/InvalidUser.html")
	}

	user, err := getUserById(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el usuario.")
		return c.RenderTemplate("App/Error.html")
	}

      	if (c.Request.Method == "GET") {
		return c.Render(user)
	}
	
	user.Name = name
	user.Lastname = lastname
	user.Role = role
	user.Email = email
	user.Created = now()

	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/user/edit/%d", id)
	}

	saveUser(&user)
	if err != nil {
		c.Flash.Error("Error modificando el usuario en la base de datos")
		return c.Redirect("/user/edit/%d", id)
	}
	c.Flash.Success("Usuario modificado")
	return c.Redirect("/user/view/%d", user.Id)
}

func (c User) View(id int) revel.Result {
	if (c.getMyRole() != "ADMIN") && (id != c.getMyUID()) {
		c.Flash.Error("No puede ver otros usuarios.")
		return c.RenderTemplate("App/InvalidUser.html")
	}
	user, err := getUserById(id)
	if err != nil {
		c.Flash.Error("Error obteniendo el usuario.")
		return c.RenderTemplate("App/Error.html")
	}
	domains, _ := domainsAssignedTo(id)
	return c.Render(user, domains)
}

func (c User) Delete(id int) revel.Result {
	// Un usuario no puede borrarse a si mismo
	// Esto implica que no se va a poder borrar siendo el ultimo usuario ADMIN del sistema.
	
	if id == c.getMyUID() {
		c.Flash.Error("Un usuario no puede borrarse a sí mismo.")
		return c.RenderTemplate("App/InvalidUser.html")
	}

	if !userExists(id) {
		c.Flash.Error("El usuario no existe.")
		return c.RenderTemplate("App/InvalidUser.html")
	}

	if hasDomains(id) {
		c.Flash.Error("No se puede borrar un usuario que administra dominios. Desvincularlo primero del dominio.")
		return c.RenderTemplate("App/Error.html")
	}

	err := deleteUserById(id)
	if err != nil {
		c.Flash.Error("Hubo un error al borrar el usuario de la base de datos.")
		return c.RenderTemplate("App/Error.html")
	}
	return c.Render()
}

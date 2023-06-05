package controllers
import (
	"github.com/revel/revel"
	"uzantoj/app/models"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type App struct {
	*revel.Controller
}

func (c App) getMyUID() int {
	tmp, _ := c.Session.Get("uid")
	myUid := int(tmp.(float64))
	return myUid
}

func (c App) getMyRole() string {
	return  c.Session["role"].(string)
}

func check_user_logged_in(c *revel.Controller) revel.Result {
	// Para lo unico que no necesita estar logueado es para:
	// - Hacer login
	// - Instalar el sistema
	// - Hacer logout
	if (c.Request.URL.String() != "/login") && (c.Request.URL.String() != "/install") && (c.Request.URL.String() != "/logout") { 
		tmp, _ := c.Session.Get("username")
		if  tmp != nil {
			// Hizo login
			return nil
		}
		return c.Redirect("/login")
	}
	return nil
}


func is_authorized_USER(c *revel.Controller) bool {
	if c.Name == "User"   && c.MethodName == "Delete"      { return false }
	if c.Name == "User"   && c.MethodName == "Add"         { return false }
	if c.Name == "User"   && c.MethodName == "List"        { return false }
	if c.Name == "User"   && c.MethodName == "SystemUsers" { return false }
	
	if c.Name == "Domain" && c.MethodName == "AddAdmin"    { return false }
	if c.Name == "Domain" && c.MethodName == "DeleteAdmin" { return false }
	if c.Name == "Domain" && c.MethodName == "Delete"      { return false }

	// Para MailUser esta habilitado para todos los methods.
	
	return true
}

func is_authorized(c *revel.Controller) revel.Result {

	// [1]
	if c.Session["role"] != nil {
		myRole :=  c.Session["role"].(string)

		if myRole == "ADMIN" {
			return nil
		} else {
			if is_authorized_USER(c) {
				return nil
			}
		}
		c.Flash.Error("El usuario no está actualizado.")
		return c.RenderTemplate("App/NotAuthorized.html")
	} else {
		return nil
	}
}
// [1]
// NOTA: Todavia no se logueo al sistema, asi que no tiene "role".
// Sin el if != nil da error si la sesion no tiene "role".


func init() {
	revel.InterceptFunc(check_user_logged_in, revel.BEFORE, &App{})
	revel.InterceptFunc(check_user_logged_in, revel.BEFORE, &Domain{})
	revel.InterceptFunc(check_user_logged_in, revel.BEFORE, &User{})

	revel.InterceptFunc(is_authorized, revel.BEFORE, &App{})
	revel.InterceptFunc(is_authorized, revel.BEFORE, &Domain{})
	revel.InterceptFunc(is_authorized, revel.BEFORE, &User{})
}


func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) InstallComplete() revel.Result {
	return c.Render()
}

// Se considera que no esta instalada si:
//   - System.installed es false
//   - No hay usuarios
//   - No hay LastGID ni LASTUID
func (c App) isInstalled() bool {
	system := &models.System{}
	DBMAP := getDB()
	defer closeDB(DBMAP)

	cant_users, err := DBMAP.SelectInt("select count(*) from users")
	if err != nil {
		revel.AppLog.Error("App.isInstalled - Error - ", err)
	}

	err = DBMAP.SelectOne(system, "select * from system")
	if err != nil {
		revel.AppLog.Error("App.isInstalled - Error - ", err)
	}

	return ((cant_users > 0) && (system.Installed) && (system.LastGID > 0) && (system.LastUID > 0))
}

func (c App) Install(username, password1, password2, lastgid, lastuid string) revel.Result {

	if c.isInstalled() {
		c.Flash.Error("La aplicación ya está instalada.")
		return c.RenderTemplate("App/AlreadyInstalled.html")
	}
	
        if (c.Request.Method == "GET") {
                return c.Render()
        }
	
	int_lastuid, _ := strconv.Atoi(lastuid)
	int_lastgid, _ := strconv.Atoi(lastgid)
	
	// TODO: Agregar validacion del lado del cliente tambien.
	c.Validation.Required(username).Message("Ingresar nombre de usuario")
	c.Validation.MaxSize(username, 15).Message("El nombre de usuario debe tener como máximo 15 caracteres")
	c.Validation.MinSize(username, 4).Message("El nombre de usuario debe tener como mínimo 4 caracteres")
	c.Validation.Match(username, regexp.MustCompile("^\\w*$")).Message("El nombre de usuario sólo puede tener letras")

	c.Validation.Required(password1).Message("Ingresar contraseña")
	c.Validation.Required(password1 == password1).Message("Las contraseñas no coinciden")

	c.Validation.Required(int_lastgid).Message("Debe indicar un GID")
	c.Validation.Range(int_lastgid, 1000, 30000).Message("El GID debe estar en el rango 1000, 30000")
	c.Validation.Range(int_lastuid, 1000, 30000).Message("El UID debe estar en el rango 1000, 30000")
	c.Validation.Required(int_lastuid).Message("Debe indicar un UID")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Install)
	}
	
	// Si llegamos aca es porque todos los parametros son correctos.
	
	system := &models.System{
		LastUID: int_lastuid,
		LastGID: int_lastgid,
		Installed: true,
		Version: 1,
	}
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	user := &models.User{
		Created: now(),
		Name: "Nombre",
		Lastname: "Apellido",
		Username: username,
		HashedPassword: bcryptPassword,
		Role: "ADMIN",
		Email: "change@me.com",		
	}
	
	DBMAP := getDB()
	defer closeDB(DBMAP)
	
	err := DBMAP.Insert(user)
	if err != nil {
		c.Flash.Error("No se pudo agregar el usuario a la base de datos.")
		revel.AppLog.Error("App.Install: Error agregando User", err)
	}

	err = DBMAP.Insert(system)
	if err != nil {
		c.Flash.Error("No se pudo agregar la información del sistema a la base de datos.")
		revel.AppLog.Error("App.Install: Error agregando System", err)
	}

	// Evitamos problemas con el F5
	return c.Redirect(App.InstallComplete)
}


func (c App) Help() revel.Result {
	return c.Render()
}


// TODO: Terminar de implementar App.Contact
func (c App) Contact() revel.Result {
	var user models.User
	contactName := revel.Config.StringDefault("contact.name", "Add a contact name")
	contactMail := revel.Config.StringDefault("contact.mail", "Add a contact mail")
	
	DBMAP := getDB()
	defer closeDB(DBMAP)

	myId, _ := c.Session.Get("uid")
	err := DBMAP.SelectOne(&user, "select * from users where id=?", myId)
	if err != nil {
		c.Flash.Error("Error al enviar la información de contacto.")
		revel.AppLog.Error("App.Contact - Error - ", err)
		return c.RenderTemplate("App/Error.html")
	}
	return c.Render(user, contactName, contactMail)

}

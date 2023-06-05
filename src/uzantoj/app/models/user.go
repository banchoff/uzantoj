package models

import (
	"fmt"
	"github.com/revel/revel"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id              int64      `db:"id"`
	Name            string     `db:"name,size:255"`
	Lastname        string     `db:"lastname,size:255"`
	Username        string     `db:"username,size:255"`
	Password        string     `db:"-"`
	HashedPassword  []byte     `db:"password,size:255"`
	Role            string     `db:"role,size:20"`
	Email           string     `db:"email,size:255"`
	Created         string	   `db:"created,size:255"`
}


func now() string {
	current_time := time.Now()
	tmp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		current_time.Year(), current_time.Month(), current_time.Day(),
		current_time.Hour(), current_time.Minute(), current_time.Second())
	return tmp
}

func UserCreate(name, lastname, username, password, role, email string) *User {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &User{
		Created: now(),
		Name: name,
		Lastname: lastname,
		Username: username,
		HashedPassword: bcryptPassword,
		Role: role,
		Email: email,	
	}
	return user
}


func (u *User) String() string {
	return fmt.Sprintf("%s %s (%s) - %s", u.Name, u.Lastname, u.Username, u.Email)
}

func (u *User) ChangePassword(newPassword string) {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	u.HashedPassword = bcryptPassword
}

func (u *User) Validate(v *revel.Validation) {
	v.Required(u.Name).Message("Debe indicar un nombre.")
	v.Match(u.Name, regexp.MustCompile("^[\\w\\s]*$")).Message("El nombre contiene caracteres inválidos.")
	v.MinSize(u.Name, 3).Message("El nombre es muy corto. Debe tener al menos 3 caracteres.")
	v.MaxSize(u.Name, 50).Message("El nombre es muy largo. Debe tener como máximo 50 caracteres.")

	v.Required(u.Lastname).Message("Debe indicar un apellido.")
	v.Match(u.Lastname, regexp.MustCompile("^[\\w\\s]*$")).Message("El apellido contiene caracteres inválidos.")
	v.MinSize(u.Lastname, 3).Message("El apellido es muy corto. Debe tener al menos 3 caracteres.")
	v.MaxSize(u.Lastname, 50).Message("El apellido es muy largo. Debe tener como máximo 50 caracteres.")
	
	v.Required(u.Email).Message("Debe indicar un email.")
	v.Email(u.Email).Message("El email dado no es válido.")

	v.Required(u.Role).Message("Debe indicar un rol para el usuario.")
	v.Match(u.Role, regexp.MustCompile("(^USER$|^ADMIN$)")).Message("El rol ingresado no es válido.")	
}

func comparePasswords(givenPassword string, userPassword []byte) bool{
	if err := bcrypt.CompareHashAndPassword(userPassword, []byte(givenPassword)); err != nil {
		return false
	}
	return true
}

func (u *User) ValidatePasswords(v *revel.Validation, aPassword string) {
	passwords_match := comparePasswords(aPassword, u.HashedPassword)
	v.Required(passwords_match).Message("Nombre de usuario o contraseña incorrectos.")	
}

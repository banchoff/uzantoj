package tests

import (
	"github.com/revel/revel/testing"
	"net/url"
	"strings"
)

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	data := url.Values{}
	data.Set("username", "TEST-USER")
	data.Set("password", "TEST-PASSWORD")
	t.PostForm("/login", data)
	println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}


func (t *AppTest) TestAddUser() {
	// Creo el User
	data := url.Values{}
	data.Set("name", "Test_Save_User_Name")
	data.Set("lastname", "Test_Save_User_Name")
	data.Set("username", "testuser")
	data.Set("password", "testing123")
	data.Set("password2", "testing123")
	data.Set("email", "test@example.com")
	data.Set("role", "USER")
	
	t.PostForm("/user/add", data)
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")

	pageContent := string(t.ResponseBody)
	t.Assert(strings.Contains(pageContent, "Usuario agregado"))

}

func (t *AppTest) After() {
	println("Tear down")
}


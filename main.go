package main

import (
	"fmt"
	"html/template"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type users struct {
	Name     string
	Username string
	password string
}

var DBusers = make(map[string]users)
var DBsession = make(map[string]string)

type errors struct {
	EmailError    string
	NameError     string
	PasswordError string
}

var errorV errors

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*"))
	DBusers["irin@gmail.com"] = users{"Irin", "irin@gmail.com", "12345"}

}

func main() {
	http.HandleFunc("/", login)
	http.HandleFunc("/home", home)

	http.HandleFunc("/logout", logout)
	fmt.Println("The Server Is Running On 8080 Port")
	http.ListenAndServe(":8080", nil)

}

// function login
func login(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	cookies, err := req.Cookie("session")
	if err == nil {
		if _, ok := DBsession[cookies.Value]; ok {
			http.Redirect(w, req, "/home", http.StatusSeeOther)
		}
	}

	if req.Method == "POST" {
		username := req.FormValue("username")
		password := req.FormValue("password")

		if _, ok := DBusers[username]; !ok {
			errorV.EmailError = "email error"
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}

		if password != DBusers[username].password {
			errorV.PasswordError = "password error"
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}
		errorV.EmailError = ""
		errorV.PasswordError = ""
		if password == DBusers[username].password {
			uid := uuid.NewV4()
			cookie := &http.Cookie{
				Name:  "session",
				Value: uid.String(),
			}
			http.SetCookie(w, cookie)
			DBsession[cookie.Value] = username
			http.Redirect(w, req, "/home", http.StatusSeeOther)
			return
		}
	}
	tmpl.ExecuteTemplate(w, "Login.html", errorV)
}

// function home
func home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	cookie, err := req.Cookie("session")

	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorV.EmailError = ""
		errorV.PasswordError = ""
		return
	}

	if _, ok := DBsession[cookie.Value]; !ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorV.EmailError = ""
		errorV.PasswordError = ""
		return
	}

	var un string
	var user users
	un = DBsession[cookie.Value]
	user = DBusers[un]

	tmpl.ExecuteTemplate(w, "home.html", user)
}

// function logout
func logout(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")

	cookie, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorV.EmailError = ""
		errorV.PasswordError = ""
		return
	}

	DBsession[cookie.Value] = ""
	delete(DBsession, cookie.Value)

	cookie = &http.Cookie{
		Name:  "session",
		Value: "",
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

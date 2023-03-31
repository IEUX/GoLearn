package auth

import (
	"Golearn/modules/client"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("M4st3r_0f_Pupp3ts")

type Creds struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Login(res http.ResponseWriter, req *http.Request) {
	//TODO
}

func Register(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		client.CORSManager(res, req)
		var user Creds
		user.Email = req.FormValue("email")
		user.Username = req.FormValue("username")
		user.Password = req.FormValue("password")
		if user.Email == "" || user.Username == "" || user.Password == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func Logout(res http.ResponseWriter, req *http.Request) {
	client.CORSManager(res, req)
	clearCookie := &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, clearCookie)
	http.Redirect(res, req, "/", http.StatusFound)
}

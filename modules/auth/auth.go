package auth

import (
	"Golearn/modules/database"
	"net/http"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("M4st3r_0f_Pupp3ts")

type Creds struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Error struct {
	Message string
}

func Login(res http.ResponseWriter, req *http.Request) {
	CORSManager(res, req)
	if req.Method == "POST" {
		var user Creds
		user.Email = req.FormValue("email")
		user.Password = req.FormValue("password")
		if user.Email == "" || user.Password == "" {
			tmpl := template.Must(template.ParseFiles("./CLIENT/static/login.gohtml"))
			_ = tmpl.Execute(res, Error{Message: "Email or password is empty"})
			return
		}
		if !database.CheckUserExist(user.Email) {
			tmpl := template.Must(template.ParseFiles("./CLIENT/static/login.gohtml"))
			_ = tmpl.Execute(res, Error{Message: "Wrong Email or Password"})
			return
		}
		dbPassword := database.GetUserPassword(user.Email)
		if dbPassword != database.HashPassword(user.Password) {
			tmpl := template.Must(template.ParseFiles("./CLIENT/static/login.gohtml"))
			_ = tmpl.Execute(res, Error{Message: "Wrong Email or Password"})
			return
		}
		expirationTime := time.Now().Add(20 * time.Minute)
		claims := &Claims{
			Email: user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JwtKey)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Error while signing token"))
			return
		}
		http.SetCookie(res, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		http.Redirect(res, req, "/", http.StatusFound)
	} else if req.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./CLIENT/static/login.gohtml"))
		_ = tmpl.Execute(res, Error{Message: ""})
	}
}

func Register(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		CORSManager(res, req)
		var user Creds
		user.Email = req.FormValue("email")
		user.Username = req.FormValue("username")
		user.Password = req.FormValue("password")
		if user.Email == "" || user.Username == "" || user.Password == "" {
			tmpl := template.Must(template.ParseFiles("./CLIENT/static/signup.gohtml"))
			_ = tmpl.Execute(res, Error{Message: "Email or password is empty"})
			return
		}
		if database.CheckUserExist(user.Email) {
			tmpl := template.Must(template.ParseFiles("./CLIENT/static/signup.gohtml"))
			_ = tmpl.Execute(res, Error{Message: "User already exist"})
			return
		}
		database.InsertUser(user.Username, user.Email, user.Password)
		Login(res, req)
	} else if req.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./CLIENT/static/signup.gohtml"))
		_ = tmpl.Execute(res, Error{Message: ""})
	}
}

func Logout(res http.ResponseWriter, req *http.Request) {
	CORSManager(res, req)
	clearCookie := &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, clearCookie)
	http.Redirect(res, req, "/", http.StatusFound)
}

func ExtractClaims(w http.ResponseWriter, r *http.Request) (database.User, bool) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return database.User{}, false
		}
		w.WriteHeader(http.StatusBadRequest)
		return database.User{}, false
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return database.User{}, false
		}
		w.WriteHeader(http.StatusBadRequest)
		return database.User{}, false
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return database.User{}, false
	}
	var user database.User
	user = database.GetUserByMail(claims.Email)
	return user, true
}

func CORSManager(res http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	res.Header().Set("Access-Control-Allow-Origin", origin)
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	res.Header().Set("Access-Control-Allow-Credentials", "true")
}

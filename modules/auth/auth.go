package auth

import (
	"Golearn/modules/client"
	"Golearn/modules/database"
	"fmt"
	"net/http"
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

func Login(res http.ResponseWriter, req *http.Request) {
	client.CORSManager(res, req)
	if req.Method == "POST" {
		var user Creds
		user.Email = req.FormValue("email")
		user.Password = req.FormValue("password")
		if user.Email == "" || user.Password == "" {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Email or password is empty"))
			return
		}
		if !database.CheckUserExist(user.Email) {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte("User does not exist"))
			return
		}
		dbPassword := database.GetUserPassword(user.Email)
		if dbPassword != database.HashPassword(user.Password) {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte("Wrong password"))
			return
		}
		expirationTime := time.Now().Add(5 * time.Minute)
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
	}
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
			res.Write([]byte("Email, username or password is empty"))
			return
		}
		if database.CheckUserExist(user.Email) {
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte("User already exist during insertion"))
			return
		}
		fmt.Println("Continue")
		database.InsertUser(user.Username, user.Email, user.Password)
		Login(res, req)
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

func ExtractClaims(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("%s\n", claims.Email)
	var user database.User
	user = database.GetUserByMail(claims.Email)
	fmt.Println(user)
}

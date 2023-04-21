package client

import (
	"fmt"
	"net/http"
)

func HomePage(res http.ResponseWriter, req *http.Request) {
	CORSManager(res, req)
	fmt.Println(req.Cookies())
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Home page"))
}

func CORSManager(res http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	res.Header().Set("Access-Control-Allow-Origin", origin)
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	res.Header().Set("Access-Control-Allow-Credentials", "true")
}

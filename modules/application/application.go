package application

import (
	"Golearn/modules/auth"
	"net/http"
)

func ReciveCode(res http.ResponseWriter, req *http.Request) {
	auth.ExtractClaims(res, req)

	//res.WriteHeader(http.StatusOK)
	//res.Write([]byte("Code recived"))

}

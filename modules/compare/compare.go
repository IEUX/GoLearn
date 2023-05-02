package compare

import (
	"Golearn/modules/auth"
	"Golearn/modules/database"
	"fmt"
	"net/http"
	"strings"
)

func GetSolution(req *http.Request) string {
	fmt.Println(req.URL.Path)
	title := strings.Split(req.URL.Path, "/exercice/")[1]
	var solution string
	err := database.GetDbInstance().QueryRow("SELECT solution FROM Exercise inner join Solution S on Exercise.ID_Exercise = S.ID_Exercise WHERE title = ?", title).Scan(&solution)
	if err != nil {
		fmt.Println(err)
	}
	return solution
}

func Compar(solution string, userResponse string, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Solution: " + solution)
	fmt.Println("User: " + userResponse)
	user, _ := auth.ExtractClaims(w, r)
	if solution == userResponse {
		err := database.GetDbInstance().QueryRow("UPDATE User SET Progression = Progression + 1 WHERE User.Username = ?", user.Name)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}

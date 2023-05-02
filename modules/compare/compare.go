package compare

import (
	"Golearn/modules/auth"
	"Golearn/modules/database"
	"fmt"
	"net/http"
)

func GetSolution(title string) string {
	var solution string
	err := database.GetDbInstance().QueryRow("SELECT solution FROM Exercise inner join Solution S on Exercise.ID_Exercise = S.ID_Exercise WHERE Title = ?", title).Scan(&solution)
	if err != nil {
		fmt.Println(err)
	}
	return solution
}

func Compar(solution string, userResponse string, w http.ResponseWriter, r *http.Request) bool {
	user, _ := auth.ExtractClaims(w, r)
	if solution == userResponse {
		database.GetDbInstance().Exec("UPDATE User SET Progression = Progression + 1 WHERE User.Username = ?", user.Name)
		return true
	} else {
		return false
	}
}

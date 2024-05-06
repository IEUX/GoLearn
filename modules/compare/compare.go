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

func GetTest(title string) string {
	var test string
	err := database.GetDbInstance().QueryRow("SELECT Test FROM Solution inner join Exercise on Solution.ID_Exercise = Exercise.ID_Exercise WHERE Title = ?", title).Scan(&test)
	if err != nil {
		fmt.Println(err)
	}
	return test
}

func Compar(exercice string, userResponse string, w http.ResponseWriter, r *http.Request) bool {
	user, _ := auth.ExtractClaims(w, r)
	if GetSolution(exercice) == userResponse {
		if database.GetExerciseByName(exercice).IdExercise-1 == user.Progression {
			database.GetDbInstance().Exec("UPDATE User SET Progression = Progression + 1 WHERE User.Username = ?", user.Name)
		}
		return true
	} else {
		return false
	}
}

package compare

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

func GetSolution() []byte {
	nbExercice := 1 //TODO: get nb exercice from database
	solutionName := "solution_" + strconv.Itoa(nbExercice)
	solution, err := os.ReadFile("./ASSETS/Solutions/" + solutionName + ".txt")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	return solution
}

func Compar(solution []byte, user []byte) {
	fmt.Println("Solution: " + string(solution))
	fmt.Println("User: " + string(user))
	if bytes.Equal(solution, user) {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}

package container

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func CreateCodeFile(user string, code string) string {
	//Create User Folder
	userFolder, err := os.MkdirTemp("./__SANDBOX/", user)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(userFolder, 0777)
	if err != nil {
		log.Fatal(err)
	}
	//COPY DOCKERFILE
	CopyDockerfile(userFolder + "/")
	//COPY GO MOD
	CopyGoMod(userFolder + "/")
	//Create Code File
	codeFile, err := os.Create(userFolder + "/main.go")
	if err != nil {
		log.Fatal(err)
	}
	defer codeFile.Close()
	//Write Code to File
	_, err = codeFile.WriteString(code)
	if err != nil {
		log.Fatal(err)
	}
	return userFolder
}

func CopyDockerfile(path string) {
	//Copy Dockerfile
	dockerFile, err := os.Open("./ASSETS/DockerfileTestCode/Dockerfile")
	if err != nil {
		log.Fatal(err)
	}
	defer dockerFile.Close()

	// Create new file
	newFile, err := os.Create(path + "Dockerfile")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, dockerFile)
	if err != nil {
		log.Fatal(err)
	}
}

func CopyGoMod(path string) {
	goMod, err := os.Open("./ASSETS/goMod/go.mod")
	if err != nil {
		log.Fatal(err)
	}
	defer goMod.Close()

	// Create new file
	newFile, err := os.Create(path + "go.mod")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, goMod)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCode(user string) []byte {
	//CLI Commmands
	build := "cd " + user[2:] + " && sudo docker build -t 'golearnbox' ."
	runDocker := "sudo docker run golearnbox"
	//-[BUILD DOCKER IMAGE]-
	cmd := exec.Command("/bin/sh", "-c", build)
	//Build errors handling
	var CompilationErr bytes.Buffer
	cmd.Stderr = &CompilationErr

	err := cmd.Run()
	if err != nil {
		var sendBackErr [][]byte
		re := regexp.MustCompile(`-{6,}\n([\s\S]*?)\n-{6,}`)
		matches := re.FindStringSubmatch(CompilationErr.String())
		if len(matches) > 1 {
			lines := bytes.Split([]byte(matches[1]), []byte("\n"))
			sendBackErr = lines[1:]
		}
		out := bytes.Join(sendBackErr, []byte("<br>"))
		return out
	}
	//-[RUN DOCKER IMAGE]-
	out, err := exec.Command("/bin/sh", "-c", runDocker).Output()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return out
}

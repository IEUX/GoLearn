package sandboxcontainer

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

func TestCode() {
	//CLI Commmands
	build := "cd ../helloWorld/ && sudo docker build -t 'golearnbox' ."
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
		fmt.Println("--- Error during compilation ---")
		fmt.Println(string(bytes.Join(sendBackErr, []byte("\n"))))

		return
	}
	//-[RUN DOCKER IMAGE]-
	result, err := cmd.Output()
	fmt.Printf("%s", result)
	out, err := exec.Command("/bin/sh", "-c", runDocker).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", out)
}

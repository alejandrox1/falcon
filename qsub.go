// qsub.go is designed to stand in for an actual qsub binary.
//
//  qsub -N P97f7c4234bfa71 -pe smp 8 -q your_queue -V -cwd -o stdout -e stderr \
// -S /bin/bash /ecoli_test/mypwatcher/wrappers/run-P97f7c4234bfa71.bash
//
package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	jobName      string
	parallelEnv  []string
	queueName    string
	exportEnv    bool
	currentDir   bool
	stdoutStream string
	stderrStream string
	jobShell     []string
)

func parseCmdArgs() {
	fmt.Println(os.Args[:])
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-N":
			if i+1 >= len(os.Args) {
				fmt.Println("Error: correct format -- ./cmd -N <string>")
				return
			}
			jobName = os.Args[i+1]
			i += 1

		case "-pe":
			if i+2 >= len(os.Args) {
				fmt.Println("Error: correct format -- ./cmd -pe <str> <str>")
				return
			}
			for _, e := range os.Args[i+1:] {
				if strings.HasPrefix(e, "-") {
					break
				}
				parallelEnv = append(parallelEnv, e)
				i += 1
			}

		case "-q":
			if i+1 >= len(os.Args) {
				fmt.Println("Error: correct format -- ./cmd -q <string>")
				return
			}
			queueName = os.Args[i+1]
			i += 1

		case "-V":
			exportEnv = true

		case "-cwd":
			currentDir = true

		case "-o":
			if i+1 >= len(os.Args) {
				fmt.Println("Error: correct format -- ./cmd -o <string>")
				return
			}
			stdoutStream = os.Args[i+1]
			i += 1

		case "-e":
			if i+1 >= len(os.Args) {
				fmt.Println("Error: correct format -- ./cmd -e <string>")
				return
			}
			stderrStream = os.Args[i+1]
			i += 1

		case "-S":
			if i+1 >= len(os.Args) {
				fmt.Println("Error: correct format -- ./cmd -S <string>...")
				return
			}
			for _, c := range os.Args[i+1:] {
				if strings.HasPrefix(c, "-") {
					break
				}
				jobShell = append(jobShell, c)
				i += 1
			}

		default:
			fmt.Println("def: ", arg)
			return
		}
	}
}

func main() {
	parseCmdArgs()

	fmt.Println("job name: ", jobName)
	fmt.Println("parallel env: ", parallelEnv)
	fmt.Println("queue: ", queueName)
	fmt.Println("export: ", exportEnv)
	fmt.Println("current directory: ", currentDir)
	fmt.Println("stdout: ", stdoutStream)
	fmt.Println("stderr: ", stderrStream)
	fmt.Println("shell: ", jobShell)
}

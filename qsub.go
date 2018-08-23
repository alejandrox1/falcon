// qsub.go is designed to stand in for an actual qsub binary.
//
//  qsub -N P97f7c4234bfa71 -pe smp 8 -q your_queue -V -cwd -o stdout -e stderr \
// -S /bin/bash /ecoli_test/mypwatcher/wrappers/run-P97f7c4234bfa71.bash
//
package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	logger *log.Logger

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
	logger.Println(os.Args[:])
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-N":
			if i+1 >= len(os.Args) {
				logger.Println("Error: correct format -- ./cmd -N <string>")
				return
			}
			jobName = os.Args[i+1]
			i += 1

		case "-pe":
			if i+2 >= len(os.Args) {
				logger.Println("Error: correct format -- ./cmd -pe <str> <str>")
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
				logger.Println("Error: correct format -- ./cmd -q <string>")
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
				logger.Println("Error: correct format -- ./cmd -o <string>")
				return
			}
			stdoutStream = os.Args[i+1]
			i += 1

		case "-e":
			if i+1 >= len(os.Args) {
				logger.Println("Error: correct format -- ./cmd -e <string>")
				return
			}
			stderrStream = os.Args[i+1]
			i += 1

		case "-S":
			if i+1 >= len(os.Args) {
				logger.Println("Error: correct format -- ./cmd -S <string>...")
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
			logger.Println("Unknown command line option: ", arg)
			os.Exit(1)
		}
	}
}

func createLogger(filename string) *log.Logger {
	prefix := "go-qsub: "
	flags := log.Lshortfile | log.Ldate

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	w := io.MultiWriter(os.Stdout, file)
	return log.New(w, prefix, flags)
}

func init() {
	logger = createLogger("qsub.log")
}

func main() {
	parseCmdArgs()

	logger.Println("job name: ", jobName)
	logger.Println("parallel env: ", parallelEnv)
	logger.Println("queue: ", queueName)
	logger.Println("export: ", exportEnv)
	logger.Println("current directory: ", currentDir)
	logger.Println("stdout: ", stdoutStream)
	logger.Println("stderr: ", stderrStream)
	logger.Println("shell: ", jobShell)

	out, err := exec.Command(jobShell[0], jobShell[1:]...).Output()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Output: %s\n", out)
	logger.Println("COMPLETED!")
}

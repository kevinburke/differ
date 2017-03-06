package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const usage = `differ [utility [argument ...]]

Execute utility with the given arguments. Then exit with an error if git reports
there are untracked changes.
`

func init() {
	flag.Usage = func() {
		os.Stderr.WriteString(usage)
	}
}

func main() {
	flag.Parse()
	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(2)
	}
	var cmd *exec.Cmd
	if len(os.Args) == 2 {
		cmd = exec.Command(os.Args[1])
	} else {
		cmd = exec.Command(os.Args[1], os.Args[2:]...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Stderr.WriteString("\n\nthe run command exited with an error; bailing")
		// actually really difficult to pass through the return code from Run so
		// just do 2
		os.Exit(2)
	}
	gitCmd := exec.Command("git", "status", "--porcelain")
	buf := new(bytes.Buffer)
	gitCmd.Stdout = buf
	gitCmd.Stderr = buf
	if err := gitCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndiffer: Error running git status --porcelain: %v\n\nOutput: %s",
			err, buf.String())
		os.Exit(2)
	}
	if buf.Len() > 0 {
		fmt.Fprintf(os.Stderr, `
Untracked or modified files present after running '%s':

%s
The command should not generate a diff. Please fix the problem and try again.
`, strings.Join(os.Args[1:], " "), buf.String())
		os.Exit(2)
	}
}

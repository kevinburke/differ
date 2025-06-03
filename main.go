package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const usage = `differ [utility [argument ...]]

Execute utility with the given arguments. Then exit with an error if git reports
there are untracked changes.
`

const Version = "1.2"

func init() {
	flag.Usage = func() {
		os.Stderr.WriteString(usage)
	}
}

func getGitDiff(ctx context.Context) string {
	diffBuf := new(bytes.Buffer)
	diffCmd := exec.CommandContext(ctx, "git", "diff", "--no-color")
	diffCmd.Stdout = diffBuf
	diffCmd.Stderr = diffBuf
	if diffErr := diffCmd.Run(); diffErr != nil {
		return ""
	}
	if diffBuf.Len() == 0 {
		return ""
	}
	bs := bufio.NewScanner(diffBuf)
	diffOutput := strings.Builder{}
	for i := 0; i < 20 && bs.Scan(); i++ {
		diffOutput.Write(bs.Bytes())
		diffOutput.WriteByte('\n')
	}
	_ = bs.Err()
	return "\nFirst few lines of the git diff:\n" + diffOutput.String()
}

func run(ctx context.Context, wd string, stderr io.Writer, args []string) int {
	if wd != "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(stderr, "Error getting current working directory: %v\n", err)
			return 2
		}
		if err := os.Chdir(wd); err != nil {
			fmt.Fprintf(stderr, "Error changing working directory to %q: %v\n", wd, err)
			return 2
		}
		defer os.Chdir(pwd)
	}
	var cmd *exec.Cmd
	if len(args) == 2 {
		cmd = exec.CommandContext(ctx, args[1])
	} else {
		cmd = exec.CommandContext(ctx, args[1], args[2:]...)
	}
	// todo encapsulation broken here
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(stderr, "\n\nthe %q command exited with an error; quitting\n", args[1])
		// actually really difficult to pass through the return code from Run so
		// just do 2
		return 2
	}
	gitCmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	buf := new(bytes.Buffer)
	gitCmd.Stdout = buf
	gitCmd.Stderr = buf
	if err := gitCmd.Run(); err != nil {
		fmt.Fprintf(stderr, `
differ: Error running git status --porcelain: %v

Output: %s`, err, buf.String())
		return 2
	}
	if buf.Len() > 0 {
		diff := getGitDiff(ctx)
		fmt.Fprintf(stderr, `
Untracked or modified files present after running '%s':

%s%s
The command should not generate a diff. Please fix the problem and try again.
`, strings.Join(args[1:], " "), buf.String(), diff)
		return 2
	}
	return 0
}

func main() {
	vsn := flag.Bool("v", false, "Print the version")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(2)
	}
	if *vsn {
		fmt.Fprintf(os.Stdout, "differ version %s\n", Version)
		os.Exit(0)
	}
	code := run(ctx, "", os.Stderr, os.Args)
	os.Exit(code)
}

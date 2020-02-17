package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	binName = "octant-operator-framework"
)

func main() {
	flag.Parse()
	for _, cmd := range flag.Args() {
		switch cmd {
		case "install":
			install()
		case "restart":
			restart()
		}
	}
}

func install() {
	b := binName
	if runtime.GOOS == "windows" {
		b = binName + ".exe"
	}

	home := os.Getenv("HOME")
	pluginDir := filepath.Join(home, ".config", "octant", "plugins")
	p := filepath.Join(pluginDir, b)

	mainPath := filepath.Join("cmd", binName)
	mainPath = "." + string(filepath.Separator) + mainPath

	runCmd("go", nil, "build", "-o", p, "-v", mainPath)
}

func restart() {
	runCmd("killall", nil, binName)
}

func runCmd(command string, env map[string]string, args ...string) {
	cmd := newCmd(command, env, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Running: %s\n", cmd.String())
	_ = cmd.Run()
}
func newCmd(command string, env map[string]string, args ...string) *exec.Cmd {
	realCommand, err := exec.LookPath(command)
	if err != nil {
		log.Fatalf("unable to find command '%s'", command)
	}

	cmd := exec.Command(realCommand, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd
}

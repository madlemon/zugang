package ssh

import (
	"os"
	"os/exec"
	"strings"
)

type CommandPattern struct {
	Name string
	Args []string
}

func (pattern *CommandPattern) inject(placeholder, value string) {
	pattern.Name = strings.ReplaceAll(pattern.Name, placeholder, value)

	for i, arg := range pattern.Args {
		pattern.Args[i] = strings.ReplaceAll(arg, placeholder, value)
	}
}

var defaultSshPattern = CommandPattern{
	Name: "plink",
	Args: []string{"{USER}@{HOST}", "-pw", "{PASSWORD}", "-no-antispoof"},
}

func Open(host, user, password string) error {
	sshPattern := defaultSshPattern

	sshPattern.inject("{HOST}", host)
	sshPattern.inject("{USER}", user)
	sshPattern.inject("{PASSWORD}", password)

	sshCmd := exec.Command(sshPattern.Name, sshPattern.Args...)

	// Set standard input/output to the current terminal
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	return sshCmd.Run()
}

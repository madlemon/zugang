package ssh

import (
	"os"
	"os/exec"
	"strings"
)

const defaultPattern = "plink {USERNAME}@{HOST} -pw {PASSWORD} -no-antispoof"

func Open(host, username, password string) error {
	sshCommand := defaultPattern

	sshCommand = strings.ReplaceAll(sshCommand, "{HOST}", host)
	sshCommand = strings.ReplaceAll(sshCommand, "{USERNAME}", username)
	sshCommand = strings.ReplaceAll(sshCommand, "{PASSWORD}", password)

	commandParts := strings.Split(sshCommand, " ")
	mainCommand := commandParts[0]
	args := commandParts[1:]

	sshCmd := exec.Command(mainCommand, args...)

	// Set standard input/output to the current terminal
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	return sshCmd.Run()
}

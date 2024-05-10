package ssh

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type CommandPattern struct {
	Name           string
	Args           []string
	PasswordPrompt bool
}

func (pattern *CommandPattern) inject(placeholder, value string) {
	pattern.Name = strings.ReplaceAll(pattern.Name, placeholder, value)

	for i, arg := range pattern.Args {
		pattern.Args[i] = strings.ReplaceAll(arg, placeholder, value)
	}
}

var openSshPattern = CommandPattern{
	Name:           "ssh",
	Args:           []string{"{USER}@{HOST}"},
	PasswordPrompt: true,
}

var plinkFlagPattern = CommandPattern{
	Name: "plink",
	Args: []string{"{USER}@{HOST}", "-pw", "{PASSWORD}", "-no-antispoof"},
}

var plinkPromptPattern = CommandPattern{
	Name:           "plink",
	Args:           []string{"{USER}@{HOST}", "-no-antispoof"},
	PasswordPrompt: true,
}

var klinkPromptPattern = CommandPattern{
	Name:           "klink",
	Args:           []string{"{USER}@{HOST}", "-no-antispoof"},
	PasswordPrompt: true,
}

var puttyFlagPattern = CommandPattern{
	Name:           "putty",
	Args:           []string{"{USER}@{HOST}", "-pw", "{PASSWORD}"},
	PasswordPrompt: true,
}

func OpenWithPattern(host, user, password string) error {
	sshPattern := openSshPattern

	sshPattern.inject("{HOST}", host)
	sshPattern.inject("{USER}", user)
	sshPattern.inject("{PASSWORD}", password)

	sshCmd := exec.Command(sshPattern.Name, sshPattern.Args...)

	stdin, err := sshCmd.StdinPipe()
	if err != nil {
		fmt.Printf("Failed to get stdin pipe: %v\n", err)
		return err
	}
	defer stdin.Close()

	stdout, err := sshCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to get stdout pipe: %v\n", err)
		return err
	}
	defer stdout.Close()

	stderr, err := sshCmd.StderrPipe()
	if err != nil {
		fmt.Printf("Failed to get stderr pipe: %v\n", err)
		return err
	}
	defer stderr.Close()

	// Start the command
	err = sshCmd.Start()
	if err != nil {
		fmt.Printf("Failed to start SSH: %v\n", err)
		return err
	}

	if sshPattern.PasswordPrompt {
		// Write the password to the command's standard input
		_, err = io.WriteString(stdin, password+"\r\n")
		if err != nil {
			fmt.Printf("Failed to write password to stdin: %v\n", err)
			return err
		}
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	go io.Copy(stdin, os.Stdin)

	err = sshCmd.Wait()
	if err != nil {
		fmt.Printf("SSH command finished with error: %v\n", err)
		return err
	}

	return err
}

func OpenWithGolang(host string, user string, password string) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := run(ctx, host, user, password); err != nil {
			log.Print(err)
		}
		cancel()
	}()

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
}

func run(ctx context.Context, host string, user string, password string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostport := fmt.Sprintf("%s:%d", host, 22)
	conn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		return fmt.Errorf("cannot connect %v: %v", hostport, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %v", err)
	}
	defer session.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("terminal make raw: %s", err)
	}
	defer terminal.Restore(fd, state)

	width, height, err := terminal.GetSize(int(os.Stdin.Fd()))

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := "xterm-256color"
	if err := session.RequestPty(term, height, width, modes); err != nil {
		return fmt.Errorf("session xterm: %s", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return fmt.Errorf("session shell: %s", err)
	}

	if err := session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}
		return fmt.Errorf("ssh: %s", err)
	}
	return nil
}

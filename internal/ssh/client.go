package ssh

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// Conf represents the configuration parameters for establishing an SSH connection.
type Conf struct {
	Address             string
	PortOverride        int
	User                string
	Password            string
	HostKeyCheckEnabled bool // HostKeyCheckEnabled determines whether host key checking is enabled or not.
}

// Connect initiates an SSH connection using the provided configuration. It handles signals for termination
// (SIGTERM, SIGINT), sets up a context, and runs the SSH session. It cancels the context when the session ends
// or upon receiving a termination signal.
func Connect(conf *Conf) {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := run(ctx, conf); err != nil {
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

// startClientConnection establishes an SSH client connection to the specified host using the provided configuration.
func startClientConnection(conf *Conf) (*ssh.Client, error) {
	clientConfig := &ssh.ClientConfig{
		User: conf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.Password),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: configureHostKeyCallback(conf.HostKeyCheckEnabled),
	}

	address := configureFinalAddress(conf)
	fmt.Println("Connecting to", address, "as", conf.User)
	conn, err := ssh.Dial("tcp", address, clientConfig)
	return conn, err
}

func configureFinalAddress(conf *Conf) string {
	if conf.PortOverride != 0 {
		// override port in address with user defined port
		host := strings.Split(conf.Address, ":")[0]
		return fmt.Sprintf("%s:%d", host, conf.PortOverride)
	}
	portMatcher := regexp.MustCompile(`:\d+$`)
	if portMatcher.MatchString(conf.Address) {
		// return address if it already contains a port
		return conf.Address
	}

	hostWithDefaultPort := fmt.Sprintf("%s:%d", conf.Address, 22)
	return hostWithDefaultPort
}

// run executes an interactive SSH session using the provided configuration.
func run(ctx context.Context, conf *Conf) error {
	conn, err := startClientConnection(conf)
	if err != nil {
		return err
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

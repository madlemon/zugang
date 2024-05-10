package ssh

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"strings"
)

// configureHostKeyCallback creates an SSH host key callback function based on whether host key checking is enabled or not.
// If host key checking is enabled, it creates a callback function using the known_hosts file. Otherwise, it uses an insecure
// callback function that ignores host key verification.
func configureHostKeyCallback(hostKeyCheckEnabled bool) ssh.HostKeyCallback {
	var hostKeyCallback ssh.HostKeyCallback
	if hostKeyCheckEnabled {
		knownHostCallback, err := hostKeyCallbackFromKnownHosts()
		if err != nil {
			log.Fatal("Failed to create HostKeyCallback:", err)
		}
		hostKeyCallback = knownHostCallback
	} else {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	return hostKeyCallback
}

// hostKeyCallbackFromKnownHosts creates an ssh.HostKeyCallback based on the known_hosts file.
// It prompts the user to confirm adding new host keys if necessary.
func hostKeyCallbackFromKnownHosts() (ssh.HostKeyCallback, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to read known_hosts: ", err)
	}
	khPath := home + "/.ssh/known_hosts"
	kh, err := knownhosts.New(khPath)
	if err != nil {
		log.Fatal("Failed to read known_hosts: ", err)
	}
	cb := ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := kh(hostname, remote, key)
		if knownhosts.IsHostKeyChanged(err) {
			return fmt.Errorf("REMOTE HOST IDENTIFICATION HAS CHANGED for host %s! This may indicate a MitM attack. If you verified that this is not an issue remove the known host key to proceed.",
				hostname)
		} else if knownhosts.IsHostUnknown(err) {
			fmt.Printf("The authenticity of host '%v' can't be established.\n", hostname)
			fmt.Printf("%v public key fingerprint is %v\n", key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
			confirmed := askAddToKnownHosts()
			if confirmed == false {
				return fmt.Errorf("unkown host key not accepted")
			}
			f, fileError := os.OpenFile(khPath, os.O_APPEND|os.O_WRONLY, 0600)
			if fileError == nil {
				defer f.Close()
				fileError = knownhosts.WriteKnownHost(f, hostname, remote, key)
			}
			if fileError == nil {
				log.Printf("Added host %s to known_hosts\n", hostname)
			} else {
				log.Printf("Failed to add host %s to known_hosts: %v\n", hostname, fileError)
				return fileError
			}
			return nil
		}
		return err
	})
	return cb, nil
}

// askAddToKnownHosts prompts the user to confirm whether they want to add the fingerprint to the known_hosts file
func askAddToKnownHosts() bool {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Do you want to add the fingerprint to known_hosts and continue connecting (yes/no)? ")

		scanner.Scan()
		answer := strings.ToLower(scanner.Text())
		if answer == "yes" {
			return true
		} else if answer == "no" {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'yes' or 'no'.")
		}
	}
}

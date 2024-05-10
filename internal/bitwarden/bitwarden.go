package bitwarden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func unlock() (string, error) {
	bwUnlockCmd := exec.Command("bw", "unlock", "--raw")

	stdin, err := bwUnlockCmd.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdin.Close()
	err = answerMasterPasswordPrompt(stdin)
	if err != nil {
		return "", err
	}

	sessionKey, err := bwUnlockCmd.Output()
	if err != nil {
		return "", err
	}

	return string(sessionKey[:]), nil
}

func answerMasterPasswordPrompt(stdin io.WriteCloser) error {
	fmt.Print("? Master password: [input is hidden]")
	masterPassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()

	go func() {
		defer stdin.Close()
		_, err := stdin.Write(masterPassword)
		if err != nil {
			log.Fatal("Error providing the master password to bw prompt", err)
		}
	}()
	return nil
}

func Lock() error {
	bwLockCmd := exec.Command("bw", "lock")
	return bwLockCmd.Run()
}

func Unlock() (string, error) {
	sessionKey, err := unlock()
	if err != nil {
		return "", fmt.Errorf("error unlocking your vault: %v", err)
	}

	return sessionKey, nil
}

func Sync() error {
	bwLockCmd := exec.Command("bw", "sync")
	return bwLockCmd.Run()
}

type Item struct {
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"login"`
}

var InvalidMasterPasswordError = fmt.Errorf("invalid master password")
var ExpiredSessionError = fmt.Errorf("bw session in env expired")

func FindSSHCredentials(host, preferredUser string) (string, string, error) {
	bwListCmd := exec.Command("bw", "list", "items", "--search", "ssh://"+host, "--pretty")

	_, sessionPresent := os.LookupEnv("BW_SESSION")
	if !sessionPresent {
		stdin, err := bwListCmd.StdinPipe()
		if err != nil {
			return "", "", err
		}
		defer stdin.Close()
		err = answerMasterPasswordPrompt(stdin)
		if err != nil {
			return "", "", err
		}
	}
	var stderr bytes.Buffer
	bwListCmd.Stderr = &stderr

	output, err := bwListCmd.Output()
	if err != nil {
		return handleBWListError(stderr, err)
	}

	var items []Item

	if err := json.Unmarshal(output, &items); err != nil {
		return "", "", fmt.Errorf("error during deserialazation: %q", err)
	}

	return extractCredentials(items, host, preferredUser)
}

func authenticateWithMasterPassword(cmd *exec.Cmd) error {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()
	return answerMasterPasswordPrompt(stdin)
}

func handleBWListError(stderr bytes.Buffer, err error) (string, string, error) {
	if strings.Contains(stderr.String(), "Invalid master password") {
		return "", "", InvalidMasterPasswordError
	}
	if strings.Contains(stderr.String(), "? Master password:") {
		log.Println(stderr.String())
		return "", "", ExpiredSessionError
	}
	log.Fatalf("Command failed with error: %v\n", err)
	return "", "", nil
}

func extractCredentials(items []Item, host, preferredUser string) (string, string, error) {
	if len(items) < 1 {
		return "", "", fmt.Errorf("did not find credentials for host %s in your vaults", host)
	}

	if preferredUser != "" {
		for _, item := range items {
			if item.Login.Username == preferredUser {
				return item.Login.Username, item.Login.Password, nil
			}
		}
	}

	if len(items) == 1 {
		return items[0].Login.Username, items[0].Login.Password, nil
	}

	return "", "", fmt.Errorf("Found multiple users within your vaults\nChoose a preferred user by providing the user flag (-u or --user)")
}

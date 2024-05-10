package bitwarden

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/exec"
)

func unlock(masterPassword []byte) (string, error) {
	bwUnlockCmd := exec.Command("bw", "unlock", "--raw")
	stdin, err := bwUnlockCmd.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdin.Close()

	go func() {
		defer stdin.Close()
		_, err := stdin.Write(masterPassword)
		if err != nil {
			log.Fatal("Error providing the master password to bw prompt", err)
		}
	}()

	sessionKey, err := bwUnlockCmd.Output()
	if err != nil {
		return "", err
	}

	return string(sessionKey[:]), nil
}

func Lock() error {
	bwLockCmd := exec.Command("bw", "lock")
	return bwLockCmd.Run()
}

func Unlock() (string, error) {
	fmt.Print("? Master password: [input is hidden]")
	masterPassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()

	sessionKey, err := unlock(masterPassword)
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

func FindSSHCredentials(host, preferredUser string) (string, string, error) {
	bwListCmd := exec.Command("bw", "list", "items", "--search", "ssh://"+host, "--pretty")
	listOutput, err := bwListCmd.Output()
	if err != nil {
		return "", "", err
	}

	var items []Item

	if err := json.Unmarshal(listOutput, &items); err != nil {
		return "", "", err
	}

	var username, password string
	if len(items) < 1 {
		return "", "", fmt.Errorf("did not find credentials for host %s in your vaults", host)
	} else if preferredUser != "" {
		username = preferredUser
		for _, item := range items {
			if item.Login.Username == preferredUser {
				password = item.Login.Password
			}
		}
		if password == "" {
			return "", "", fmt.Errorf("did not find a password for preferred user %s on host %s in your vaults", preferredUser, host)
		}
	} else if len(items) == 1 {
		username = items[0].Login.Username
		password = items[0].Login.Password
	} else {
		var usernames []string
		for _, item := range items {
			usernames = append(usernames, item.Login.Username)
		}
		return "", "", fmt.Errorf("Found multiple users within your vaults: %v\nChoose a preferred a user by providing the user flag (-u or --user)", usernames)
	}

	return username, password, nil
}

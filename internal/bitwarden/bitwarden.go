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

// CheckForExecutable checks if the Bitwarden CLI executable is available in the system's PATH.
func CheckForExecutable() {
	_, err := exec.LookPath("bw")
	if err != nil {
		log.Fatal("Bitwarden CLI executable not found: ", err)
	}
}

// answerMasterPasswordPrompt prompts the user for the master password to provide to the Bitwarden CLI when prompted.
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

// Sync synchronizes local data with the Bitwarden server.
func Sync() error {
	bwLockCmd := exec.Command("bw", "sync")
	return bwLockCmd.Run()
}

// BwItem represents an item returned by Bitwarden
type bwItem struct {
	Id    string `json:"id"`
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Uris     []uri  `json:"uris"`
	} `json:"login"`
}

type uri struct {
	Uri string `json:"uri"`
}

type VaultItem struct {
	Address  string
	Username string
	Password string
}

const sshProtocolPrefix = "ssh://"

var InvalidMasterPasswordError = fmt.Errorf("invalid master password")

// FindVaultItem searches for SSH credentials associated with the specified host in Bitwarden.
func FindVaultItem(host, preferredUser string) (*VaultItem, error) {

	bwListCmd := exec.Command("bw", "list", "items", "--search", sshProtocolPrefix+host, "--pretty")

	stdin, err := bwListCmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()
	err = answerMasterPasswordPrompt(stdin)
	if err != nil {
		return nil, err
	}

	var stderr bytes.Buffer
	bwListCmd.Stderr = &stderr

	output, err := bwListCmd.Output()
	if err != nil {
		if strings.Contains(stderr.String(), "Invalid master password") {
			return nil, InvalidMasterPasswordError
		}
		log.Fatalf("Command failed with error: %v\n", err)
	}

	var items []bwItem

	if err := json.Unmarshal(output, &items); err != nil {
		return nil, fmt.Errorf("error during deserialazation: %q", err)
	}

	bwItem, err := findMatchingItem(items, host, preferredUser)
	if err != nil {
		return nil, err
	}

	for _, uri := range bwItem.Login.Uris {
		if strings.HasPrefix(uri.Uri, sshProtocolPrefix) {
			address := strings.ReplaceAll(uri.Uri, sshProtocolPrefix, "")
			vaultItem := &VaultItem{
				Username: bwItem.Login.Username,
				Password: bwItem.Login.Password,
				Address:  address,
			}
			return vaultItem, nil
		}
	}

	return nil, fmt.Errorf("cannot find matching uri")
}

func findMatchingItem(items []bwItem, host, preferredUser string) (*bwItem, error) {
	if len(items) < 1 {
		return nil, fmt.Errorf("did not find credentials for host %s in your vaults", host)
	}

	if preferredUser != "" {
		for _, item := range items {
			if item.Login.Username == preferredUser {
				return &item, nil
			}
		}
	}

	if len(items) == 1 {
		return &items[0], nil
	}

	return nil, fmt.Errorf("Found multiple users within your vaults\nChoose a preferred user by providing the user flag (-u or --user)")
}

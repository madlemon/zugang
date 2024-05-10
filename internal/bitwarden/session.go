package bitwarden

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var sessionFilePath = filepath.Join(os.TempDir(), "zugang_session")

func SaveSession(sessionKey string) error {
	err := os.WriteFile(sessionFilePath, []byte(sessionKey), 0644)
	if err != nil {
		return fmt.Errorf("failed writing to session file")
	}

	return nil
}

func LoadSession() (string, error) {
	if _, err := os.Stat(sessionFilePath); errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	data, err := os.ReadFile(sessionFilePath)
	if err != nil {
		return "", fmt.Errorf("failed reading from session file")
	}
	sessionKey := strings.TrimSpace(string(data))

	return sessionKey, nil
}

func LoadOrCreateSessionIfNotExists() error {
	sessionKey, err := LoadSession()
	if sessionKey == "" {
		sessionKey, err = Unlock()
		if err != nil {
			return fmt.Errorf("error unlocking your vault: %q", err)

		}
		err = SaveSession(sessionKey)
		if err != nil {
			return fmt.Errorf("error saving your session: %q", err)
		}
	}
	_ = os.Setenv("BW_SESSION", sessionKey[:])

	return nil
}

func DiscardSession() error {
	if _, err := os.Stat(sessionFilePath); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	err := os.Remove(sessionFilePath)
	if err != nil {
		return fmt.Errorf("error removing stored session: %q", err)
	}
	return nil
}

package ssh

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
)

// SSHClient структура для параметров подключения.
type SSHClient struct {
	Username       string
	PrivateKeyPath string
}

// GetDefaultUsername возвращает имя текущего пользователя в системе.
func GetDefaultUsername() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %v", err)
	}
	return currentUser.Username
}

// Connect выполняет SSH-подключение к указанному хосту.
func (s *SSHClient) Connect(host string) error {
	username := s.Username
	if username == "" {
		username = GetDefaultUsername()
	}

	sshArgs := []string{}
	if s.PrivateKeyPath != "" {
		sshArgs = append(sshArgs, "-i", s.PrivateKeyPath)
	}

	sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", username, host))

	fmt.Printf("Connecting to %s as %s...\n", host, username)
	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

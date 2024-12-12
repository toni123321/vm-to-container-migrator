package analyze

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func createBaseOutputDir(baseOutputDir string) {
	if _, err := os.Stat(baseOutputDir); os.IsNotExist(err) {
		os.Mkdir(baseOutputDir, 0755)
	}
}

func connect_ssh(user string, host string, port string, privateKeyPath string) *ssh.Client {
	// Create SSH client configuration
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the VM
	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	return client
}

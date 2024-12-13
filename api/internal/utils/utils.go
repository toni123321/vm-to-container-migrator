package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

func CreateBaseOutputDir(baseOutputDir string) {
	if _, err := os.Stat(baseOutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(baseOutputDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}
}

func ConnectToSsh(user string, host string, port string, privateKeyPath string) *ssh.Client {
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

func AllowSudoForTar(usr string) {
	// Check if the current user can run tar without sudo already
	cmdCheck := "sudo -l | grep tar"
	cmdCheckExec := exec.Command("sh", "-c", cmdCheck)
	if err := cmdCheckExec.Run(); err == nil {
		fmt.Println("Current user can run tar without sudo!")
		return
	}

	// Current user can run tar without sudo through visudo
	cmd := fmt.Sprintf("echo \"%s ALL=(ALL) NOPASSWD: /bin/tar\" | sudo EDITOR='tee -a' visudo", usr)
	cmdExec := exec.Command("sh", "-c", cmd)

	fmt.Println("Allowing sudo for the `tar` command ...")

	err := cmdExec.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error allowing sudo for the `tar` command: %s\n", err)
	}

	fmt.Println("Successfully allowed sudo for the `tar` command!")
}

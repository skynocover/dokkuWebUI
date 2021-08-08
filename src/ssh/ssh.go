package ssh

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

type sshClient struct {
	c *ssh.Client
}

var Client = &sshClient{}

func (client *sshClient) Init() error {
	signer, err := ssh.ParsePrivateKey([]byte(os.Getenv("PRIVATE_KEY")))
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User:            os.Getenv("SSH_USER"),
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	c, err := ssh.Dial("tcp", os.Getenv("SSH_SERVER")+":22", config)
	if err != nil {
		return err
	}
	Client.c = c

	session, err := c.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("dokku apps:list"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())

	return nil
}

func (client *sshClient) Run(args string) (string, error) {
	session, err := client.c.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(args); err != nil {
		return "", err
	}
	return b.String(), nil
}

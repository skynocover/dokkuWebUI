package ssh

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

type sshClient struct {
	hasKey bool
	c      *ssh.Client
}

var Client = &sshClient{hasKey: false}

func (client *sshClient) Init(privateKey []byte) error {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		log.Println("Failed to parse")
		return err
	}

	config := &ssh.ClientConfig{
		User:            os.Getenv("SSH_USER"),
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	c, err := ssh.Dial("tcp", os.Getenv("SSH_SERVER")+":22", config)
	if err != nil {
		log.Println("Failed to Dial")
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

	client.hasKey = true
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

func (client *sshClient) CheckSSHKey() bool {
	return client.hasKey
}

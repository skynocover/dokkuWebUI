package ssh

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

type sshClient struct {
	c *ssh.Client
}

var Client = &sshClient{}

func (client *sshClient) Init() error {
	privateKey := os.Getenv("PRIVATE_KEY")
	privateKey = parse(privateKey)

	log.Printf("PRIVATE_KEY: %s\n", privateKey)

	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
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

var (
	singleQuotesRegex  = regexp.MustCompile(`\A'(.*)'\z`)
	doubleQuotesRegex  = regexp.MustCompile(`\A"(.*)"\z`)
	escapeRegex        = regexp.MustCompile(`\\.`)
	unescapeCharsRegex = regexp.MustCompile(`\\([^$])`)
)

func parse(value string) string {
	singleQuotes := singleQuotesRegex.FindStringSubmatch(value)

	doubleQuotes := doubleQuotesRegex.FindStringSubmatch(value)

	if singleQuotes != nil || doubleQuotes != nil {
		// pull the quotes off the edges
		value = value[1 : len(value)-1]
	}

	if doubleQuotes != nil {
		// expand newlines
		value = escapeRegex.ReplaceAllStringFunc(value, func(match string) string {
			c := strings.TrimPrefix(match, `\`)
			switch c {
			case "n":
				return "\n"
			case "r":
				return "\r"
			default:
				return match
			}
		})
		// unescape characters
		value = unescapeCharsRegex.ReplaceAllString(value, "$1")
	}

	return value
}

var expandVarRegex = regexp.MustCompile(`(\\)?(\$)(\()?\{?([A-Z0-9_]+)?\}?`)

func expandVariables(v string, m map[string]string) string {
	return expandVarRegex.ReplaceAllStringFunc(v, func(s string) string {
		submatch := expandVarRegex.FindStringSubmatch(s)

		if submatch == nil {
			return s
		}
		if submatch[1] == "\\" || submatch[2] == "(" {
			return submatch[0][1:]
		} else if submatch[4] != "" {
			return m[submatch[4]]
		}
		return s
	})
}

const doubleQuoteSpecialChars = "\\\n\r\"!$`"

func doubleQuoteEscape(line string) string {
	for _, c := range doubleQuoteSpecialChars {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		line = strings.Replace(line, string(c), toReplace, -1)
	}
	return line
}

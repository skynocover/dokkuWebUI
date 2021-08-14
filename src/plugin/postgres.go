package plugin

import (
	"dokkuwebui/src/ssh"
)

type postgres struct{}

func (p *postgres) Install() error {
	_, err := ssh.Client.Run("sudo dokku plugin:install https://github.com/dokku/dokku-postgres.git postgres")
	return err
}

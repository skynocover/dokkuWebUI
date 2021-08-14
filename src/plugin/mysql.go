package plugin

import "dokkuwebui/src/ssh"

type mysql struct{}

func (m *mysql) Install() error {
	_, err := ssh.Client.Run("sudo dokku plugin:install https://github.com/dokku/dokku-mysql.git mysql")
	return err
}

package plugin

import "dokkuwebui/src/ssh"

type mariadb struct{}

func (m *mariadb) Install() error {
	_, err := ssh.Client.Run("sudo dokku plugin:install https://github.com/dokku/dokku-mariadb.git mariadb")

	return err
}

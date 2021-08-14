package plugin

import "dokkuwebui/src/ssh"

type mongo struct{}

func (m *mongo) Install() error {
	_, err := ssh.Client.Run("sudo dokku plugin:install https://github.com/dokku/dokku-mongo.git mongo")

	return err
}

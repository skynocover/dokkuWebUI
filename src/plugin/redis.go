package plugin

import "dokkuwebui/src/ssh"

type redis struct{}

func (r *redis) Install() error {
	_, err := ssh.Client.Run("sudo dokku plugin:install https://github.com/dokku/dokku-redis.git redis")

	return err
}

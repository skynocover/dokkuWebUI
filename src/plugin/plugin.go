package plugin

// var plugin = []string{"postgres", "mariadb", "mongo", "redis", "mysql"}

func Generate(name string) Plugin {
	switch name {
	case "postgres":
		return &postgres{}
	case "mariadb":
		return &mariadb{}
	case "mongo":
		return &mongo{}
	case "redis":
		return &redis{}
	case "mysql":
		return &mysql{}
	default:
		return nil
	}
}

type Plugin interface {
	Install() error
}

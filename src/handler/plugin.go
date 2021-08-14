package handler

import (
	"dokkuwebui/src/plugin"
	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
)

func Databases(ctx iris.Context) {
	r, err := ssh.Client.Run("dokku plugin:list")
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"databases": parsePlugin(r)}))
}

func DatabaseInstall(ctx iris.Context) {
	database := ctx.Params().GetString("database")

	db := plugin.Generate(database)
	if db == nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	if err := db.Install(); err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}

	ctx.Write(resp.SUCCESS.ToBytes())

}

func DatabaseEnable(ctx iris.Context) {
	database := ctx.Params().GetString("database")
	enable := ctx.Params().GetString("enable")
	switch enable {
	case "enable":
		_, err := ssh.Client.Run(fmt.Sprintf("sudo dokku plugin:enable %s", database))
		if err != nil {
			ctx.Write(resp.SystemError.ToBytesWithErr(err))
			return
		}
	case "disable":
		_, err := ssh.Client.Run(fmt.Sprintf("sudo dokku plugin:disable %s", database))
		if err != nil {
			ctx.Write(resp.SystemError.ToBytesWithErr(err))
			return
		}
	default:
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

func DatabaseUninstall(ctx iris.Context) {
	database := ctx.Params().GetString("database")
	_, err := ssh.Client.Run(fmt.Sprintf("sudo dokku plugin:uninstall %s", database))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

type tplugin struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func parsePlugin(r string) (plugins []tplugin) {
	plugins = []tplugin{}
	lists := strings.Split(r, "\n")
	for _, list := range lists {
		plug := strings.Fields(list)
		if len(plug) > 0 {
			if plugin.Generate(plug[0]) != nil {
				plugins = append(plugins, tplugin{Name: plug[0], Enabled: plug[2] == "enabled"})
			}
		}
	}
	return
}

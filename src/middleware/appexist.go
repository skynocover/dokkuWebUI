package middleware

import (
	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"fmt"
	"log"

	"github.com/kataras/iris/v12"
)

func AppExists(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")

	r, err := ssh.Client.Run(fmt.Sprintf("dokku apps:exists %s", appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	log.Println("r: ", r)
	if r != "" {
		ctx.Write(resp.ErrorAppNotExist.ToBytes())
		return
	}
	ctx.Values().Set("AppName", appName)
	ctx.Next()
}

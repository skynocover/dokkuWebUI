package handler

import (
	"fmt"

	"github.com/kataras/iris/v12"

	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"dokkuwebui/src/utils"
)

func GitInit(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	_, err := ssh.Client.Run(fmt.Sprintf("dokku git:initialize %s", appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

func GitReport(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	r, err := ssh.Client.Run(fmt.Sprintf("dokku git:report %s", appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	report := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
}

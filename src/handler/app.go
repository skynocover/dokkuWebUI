package handler

import (
	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"dokkuwebui/src/utils"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
)

func AppList(ctx iris.Context) {
	r, err := ssh.Client.Run("dokku apps:list")
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	rs := strings.Split(r, "\n")
	list := rs[1 : len(rs)-1]

	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"list": list}))
}

func AppReport(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	r, err := ssh.Client.Run(fmt.Sprintf("dokku apps:report %s", appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	report := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
}

func AppCreate(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	r, err := ssh.Client.Run(fmt.Sprintf("dokku apps:create %s", appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	if strings.HasPrefix(r, "!") {
		ctx.Write(resp.ErrorAppAlreadyExist.ToBytes())
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

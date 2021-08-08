package handler

import (
	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
)

func Logs(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	num := ctx.Params().GetString("num")
	if _, err := strconv.Atoi(num); err != nil {
		ctx.Write(resp.ErrorParameter.ToBytesWithErr(err))
		return
	}

	report, err := ssh.Client.Run(fmt.Sprintf("dokku logs -n %s %s", num, appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	// report := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
}

func LogsErr(ctx iris.Context) {
	appName := ctx.FormValue("appName")
	if appName == "" {
		report, err := ssh.Client.Run(fmt.Sprintf("dokku logs:failed --all"))
		if err != nil {
			ctx.Write(resp.SystemError.ToBytesWithErr(err))
			return
		}
		// report := utils.Parse(r)
		ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
	} else {
		report, err := ssh.Client.Run(fmt.Sprintf("dokku logs:failed %s", appName))
		if err != nil {
			ctx.Write(resp.SystemError.ToBytesWithErr(err))
			return
		}
		// report := utils.Parse(r)
		ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
	}
}

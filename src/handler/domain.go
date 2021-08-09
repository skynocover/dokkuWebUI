package handler

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"

	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"dokkuwebui/src/utils"
)

func DomainReport(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	r, err := ssh.Client.Run(fmt.Sprintf("dokku domains:report %s", appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	report := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
}

func DomainEnable(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	enable := ctx.Params().GetString("enable")
	if enable == "enable" {
		_, err := ssh.Client.Run(fmt.Sprintf("dokku domains:enable %s", appName))
		if err != nil {
			ctx.Write(resp.SystemError.ToBytesWithErr(err))
			return
		}
	} else if enable == "disable" {
		_, err := ssh.Client.Run(fmt.Sprintf("dokku domains:disable %s", appName))
		if err != nil {
			ctx.Write(resp.SystemError.ToBytesWithErr(err))
			return
		}
	} else {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

func DomainAdd(ctx iris.Context) {
	var domains = []string{}
	if err := ctx.ReadJSON(&domains); err != nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}
	appName := ctx.Params().GetString("appName")

	_, err := ssh.Client.Run(fmt.Sprintf("dokku domains:add %s %s", appName, strings.Join(domains, " ")))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

func DomainRemove(ctx iris.Context) {
	var domains = []string{}
	if err := ctx.ReadJSON(&domains); err != nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	appName := ctx.Params().GetString("appName")
	_, err := ssh.Client.Run(fmt.Sprintf("dokku domains:remove %s %s", appName, strings.Join(domains, " ")))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

/////////////////////////////// global ////////////////////////////////
func DomainGlobal(ctx iris.Context) {
	r, err := ssh.Client.Run(fmt.Sprintf("dokku domains:report --global"))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	report := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
}

func DomainGlobalSet(ctx iris.Context) {
	var domains = []string{}
	if err := ctx.ReadJSON(&domains); err != nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	_, err := ssh.Client.Run(fmt.Sprintf("dokku domains:add-global %s", strings.Join(domains, " ")))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

func DomainGlobalRemove(ctx iris.Context) {
	domain := ctx.Params().GetString("domain")
	_, err := ssh.Client.Run(fmt.Sprintf("dokku domains:remove-global %s", domain))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

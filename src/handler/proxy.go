package handler

import (
	"fmt"

	"github.com/kataras/iris/v12"

	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"dokkuwebui/src/utils"
)

func ProxyReport(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	r, err := ssh.Client.Run(fmt.Sprintf("dokku proxy:report %s", appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	report := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
}

type Proxy struct {
	Scheme        string `json:"scheme"`
	HostPort      string `json:"hostPort"`
	ContainerPort string `json:"containerPort"`
}

func ProxyAdd(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	var proxy Proxy
	if err := ctx.ReadJSON(&proxy); err != nil {
		ctx.Write(resp.ErrorParsingJSON.ToBytesWithErr(err))
		return
	}

	proxymap := fmt.Sprintf("%s:%s:%s ", proxy.Scheme, proxy.HostPort, proxy.ContainerPort)
	_, err := ssh.Client.Run(fmt.Sprintf("dokku proxy:ports-add %s %s", appName, proxymap))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

func ProxyRemove(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	var proxy Proxy
	if err := ctx.ReadJSON(&proxy); err != nil {
		ctx.Write(resp.ErrorParsingJSON.ToBytesWithErr(err))
		return
	}

	proxymap := fmt.Sprintf("%s:%s:%s ", proxy.Scheme, proxy.HostPort, proxy.ContainerPort)
	_, err := ssh.Client.Run(fmt.Sprintf("dokku proxy:ports-remove %s %s", appName, proxymap))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

func ProxyEnable(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	enable := ctx.Params().GetString("enable")
	if enable == "enable" {
		_, err := ssh.Client.Run(fmt.Sprintf("dokku proxy:enable %s", appName))
		if err != nil {
			ctx.Write(resp.SystemError.ToBytesWithErr(err))
			return
		}
	} else if enable == "disable" {
		_, err := ssh.Client.Run(fmt.Sprintf("dokku proxy:disable %s", appName))
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

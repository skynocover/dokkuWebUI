package handler

import (
	"fmt"
	"log"

	"github.com/kataras/iris/v12"

	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"dokkuwebui/src/utils"
)

type Setting struct {
	Restart bool `json:"restart"`
	Configs []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"config"`
}

func ConfigShow(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")

	r, err := ssh.Client.Run(fmt.Sprintf("dokku config:show %s ", appName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	report := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
}

func ConfigSet(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	var setting Setting
	if err := ctx.ReadJSON(&setting); err != nil {
		ctx.Write(resp.ErrorParsingJSON.ToBytesWithErr(err))
		return
	}

	config := ""
	for i := range setting.Configs {
		config += fmt.Sprintf("%s=%s ", setting.Configs[i].Key, setting.Configs[i].Value)
	}
	log.Println(config)

	var err error
	if setting.Restart {
		_, err = ssh.Client.Run(fmt.Sprintf("dokku config:set %s %s", appName, config))
	} else {
		_, err = ssh.Client.Run(fmt.Sprintf("dokku config:set --no-restart %s %s", appName, config))
	}
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}

	ctx.Write(resp.SUCCESS.ToBytes())
}

func ConfigUnset(ctx iris.Context) {
	appName := ctx.Params().GetString("appName")
	key := ctx.Params().GetString("key")
	restart := true
	rs := ctx.FormValue("restart")
	if rs == "false" {
		restart = false
	}

	var err error
	if restart {
		_, err = ssh.Client.Run(fmt.Sprintf("dokku config:unset %s %s", appName, key))
	} else {
		_, err = ssh.Client.Run(fmt.Sprintf("dokku config:unset --no-restart %s %s", appName, key))
	}
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

////////////////////////////////// global //////////////////////////////////

func ConfigGlobalShow(ctx iris.Context) {
	r, err := ssh.Client.Run("dokku config:show --global ")
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	report := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"report": report}))
}

func ConfigGlobalSet(ctx iris.Context) {
	var setting Setting
	if err := ctx.ReadJSON(&setting); err != nil {
		ctx.Write(resp.ErrorParsingJSON.ToBytesWithErr(err))
		return
	}

	config := ""
	for i := range setting.Configs {
		config += fmt.Sprintf("%s=%s ", setting.Configs[i].Key, setting.Configs[i].Value)
	}
	log.Println(config)

	var err error
	if setting.Restart {
		_, err = ssh.Client.Run(fmt.Sprintf("dokku config:set --global %s", config))
	} else {
		_, err = ssh.Client.Run(fmt.Sprintf("dokku config:set --no-restart --global %s", config))
	}
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}

	ctx.Write(resp.SUCCESS.ToBytes())
}

func ConfigGlobalUnset(ctx iris.Context) {
	key := ctx.Params().GetString("key")
	restart := true
	rs := ctx.FormValue("restart")
	if rs == "false" {
		restart = false
	}

	var err error
	if restart {
		_, err = ssh.Client.Run(fmt.Sprintf("dokku config:unset --global %s", key))
	} else {
		_, err = ssh.Client.Run(fmt.Sprintf("dokku config:unset --no-restart --global %s", key))
	}
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

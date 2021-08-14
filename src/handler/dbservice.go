package handler

import (
	"dokkuwebui/src/plugin"
	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
	"dokkuwebui/src/utils"
	"fmt"
	"log"
	"strings"

	"github.com/kataras/iris/v12"
)

func DBServices(ctx iris.Context) {
	database := ctx.Params().GetString("database")

	db := plugin.Generate(database)
	if db == nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	r, err := ssh.Client.Run(fmt.Sprintf("dokku %s:list", database))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}

	services := utils.ParseNewLine(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"services": services}))
}

func DBServiceInfo(ctx iris.Context) {
	database := ctx.Params().GetString("database")
	service := ctx.Params().GetString("service")

	db := plugin.Generate(database)
	if db == nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	r, err := ssh.Client.Run(fmt.Sprintf("dokku %s:info %s", database, service))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}

	services := utils.Parse(r)
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"info": services}))
}

func DBServiceAdd(ctx iris.Context) {
	database := ctx.Params().GetString("database")

	var params struct {
		Service      string `json:"service"`
		Password     string `json:"password"`
		RootPassword string `json:"rootPassword"`
	}

	if err := ctx.ReadJSON(&params); err != nil {
		ctx.Write(resp.ErrorParsingJSON.ToBytesWithErr(err))
		return
	}

	if db := plugin.Generate(database); db == nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	var param = ""
	if params.Password != "" {
		param = fmt.Sprintf("%s -p %s ", param, params.Password)
	}

	if params.RootPassword != "" {
		param = fmt.Sprintf("%s -r %s ", param, params.RootPassword)
	}

	command := fmt.Sprintf("dokku %s:create %s %s", database, params.Service, param)
	log.Println("command: ", command)

	_, err := ssh.Client.Run(command)
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

// func DBServiceDestroy(ctx iris.Context) {
// 	database := ctx.Params().GetString("database")
// 	service := ctx.Params().GetString("service")

// 	db := plugin.Generate(database)
// 	if db == nil {
// 		ctx.Write(resp.ErrorParameter.ToBytes())
// 		return
// 	}

// 	_, err := ssh.Client.Run(fmt.Sprintf("dokku %s:destroy %s", database, service))
// 	if err != nil {
// 		ctx.Write(resp.SystemError.ToBytesWithErr(err))
// 		return
// 	}
// 	ctx.Write(resp.SUCCESS.ToBytes())
// }

func DBStart(ctx iris.Context) {
	var params struct {
		Database string `json:"database"`
		Service  string `json:"service"`
		Start    string `json:"start"`
	}

	if err := ctx.ReadJSON(&params); err != nil {
		ctx.Write(resp.ErrorParsingJSON.ToBytesWithErr(err))
		return
	}

	db := plugin.Generate(params.Database)
	if db == nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	switch params.Start {
	case "start":
	case "restart":
	case "stop":
	default:
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	_, err := ssh.Client.Run(fmt.Sprintf("dokku %s:%s %s ", params.Database, params.Start, params.Service))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

func DBLinks(ctx iris.Context) {
	database := ctx.Params().GetString("database")
	service := ctx.Params().GetString("service")

	db := plugin.Generate(database)
	if db == nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	r, err := ssh.Client.Run(fmt.Sprintf("dokku %s:links %s", database, service))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}

	links := strings.Split(r, "\n")
	if len(links) > 0 {
		links = links[:len(links)-1]
	}
	ctx.Write(resp.SUCCESS.ToBytesWithObject(iris.Map{"links": links}))
}

func DBLink(ctx iris.Context) {
	var params struct {
		Database string `json:"database"`
		Service  string `json:"service"`
		AppName  string `json:"appName"`
		Link     bool   `json:"link"`
	}

	if err := ctx.ReadJSON(&params); err != nil {
		ctx.Write(resp.ErrorParsingJSON.ToBytesWithErr(err))
		return
	}

	db := plugin.Generate(params.Database)
	if db == nil {
		ctx.Write(resp.ErrorParameter.ToBytes())
		return
	}

	var link = "unlink"
	if params.Link {
		link = "link"
	}

	_, err := ssh.Client.Run(fmt.Sprintf("dokku %s:%s %s %s", params.Database, link, params.Service, params.AppName))
	if err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

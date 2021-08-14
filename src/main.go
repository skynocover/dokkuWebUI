package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"

	"dokkuwebui/src/handler"
	"dokkuwebui/src/middleware"
	"dokkuwebui/src/ssh"
)

func main() {

	matches, err := filepath.Glob(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		if len(matches) > 0 {
			if err := godotenv.Load(); err != nil {
				log.Fatal("Error loading .env file: ", err)
			}
			if os.Getenv("ENV") == "development" {
				dat, err := ioutil.ReadFile("../dokku")
				if err != nil {
					log.Fatal("Error read ssh key")
				}
				ssh.Client.Init(dat)
			}
		}
	}

	app := iris.New()

	log.Println(fmt.Sprintf("Serving at localhost:%s", os.Getenv("SERVER_LISTEN")))

	app.HandleDir("/", "../public")

	app.Get("/api/version", handler.Version)
	app.Post("/api/account/login", handler.Login)
	app.Post("/api/account/logout", handler.Logout)
	// ssh key
	app.Post("/api/ssh/upload", handler.Upload)

	// database
	app.Get("/api/databases", handler.Databases)
	app.Post("/api/database/{database}", handler.DatabaseInstall)
	app.Patch("/api/database/{database}/{enable}", handler.DatabaseEnable)
	app.Delete("/api/database/{database}", handler.DatabaseUninstall)

	// database service
	app.Get("/api/dbservice/{database}", handler.DBServices)
	app.Get("/api/dbservice/{database}/{service}", handler.DBServiceInfo)
	app.Post("/api/dbservice/{database}", handler.DBServiceAdd)
	// app.Delete("/api/dbservice/{database}/{service}", handler.DBServiceDestroy)
	app.Patch("/api/dbservice", handler.DBStart)

	app.Get("/api/dblinks/{database}/{service}", handler.DBLinks)
	app.Post("/api/dblink", handler.DBLink)

	api := app.Party("/api", middleware.CheckPrivateKey, handler.CheckExpire)
	{
		// redirect
		api.Get("/redirect", handler.Redirect)

		// log
		api.Get("/logs/{appName}/{num}", handler.Logs)
		// app
		api.Get("/apps", handler.AppList)
		api.Get("/app/{appName}", handler.AppReport)
		api.Post("/app/{appName}", handler.AppCreate)

		// git
		api.Get("/git/{appName}", handler.GitReport)
		api.Post("/git/{appName}", handler.GitInit)

		// proxy
		api.Get("/proxy/{appName}", handler.ProxyReport)
		api.Post("/proxy/{appName}", handler.ProxyAdd)
		api.Delete("/proxy/{appName}", handler.ProxyRemove)
		api.Patch("/proxy/{appName}/{enable}", handler.ProxyEnable)

		// domain
		api.Get("/domain/{appName}", handler.DomainReport)
		api.Post("/domain/{appName}", handler.DomainAdd)
		api.Delete("/domain/{appName}", handler.DomainRemove)
		api.Patch("/domain/{appName}/{enable}", handler.DomainEnable)
		// domain global
		api.Get("/globaldomain", handler.DomainGlobal)
		api.Post("/globaldomain", handler.DomainGlobalSet)
		api.Delete("/globaldomain/{domain}", handler.DomainGlobalRemove)

		// config
		api.Get("/config/{appName}", handler.ConfigShow)
		api.Post("/config/{appName}", handler.ConfigSet)
		api.Delete("/config/{appName}/{key}", handler.ConfigUnset)
		// config global
		api.Get("/globalconfig", handler.ConfigGlobalShow)
		api.Post("/globalconfig", handler.ConfigGlobalSet)
		api.Delete("/globalconfig/{key}", handler.ConfigGlobalUnset)
	}

	if err := app.Run(
		iris.Addr(":"+os.Getenv("SERVER_LISTEN")),
		iris.WithoutPathCorrection,
		iris.WithoutServerError(iris.ErrServerClosed),
	); err != nil {
		log.Fatal("failed run app: ", err)
	}
}

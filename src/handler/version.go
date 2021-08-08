package handler

import (
	"os"

	"github.com/kataras/iris/v12"
)

func Version(ctx iris.Context) {
	ctx.JSON(iris.Map{"Listen": os.Getenv("SERVER_LISTEN"), "Version": os.Getenv("VERSION"), "Env": os.Getenv("ENV")})
}

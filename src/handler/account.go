package handler

import (
	"os"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"

	"dokkuwebui/src/resp"
)

var (
	sessionCookie = "sesson"
	SessionConfig = sessions.Config{Cookie: sessionCookie, Expires: 30 * time.Minute}
	sess          = sessions.New(SessionConfig)
)

func Redirect(ctx iris.Context) {
	ctx.Write(resp.SUCCESS.ToBytes())
}

func CheckExpire(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		ctx.Write(resp.SessionExpired.ToBytes())
		return
	}
	ctx.Next()
}

func Login(ctx iris.Context) {
	var user struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	if err := ctx.ReadJSON(&user); err != nil {
		ctx.Write(resp.ErrorParsingJSON.ToBytesWithErr(err))
		return
	}

	if user.Account != os.Getenv("ACCOUNT") || user.Password != os.Getenv("PASSWORD") {
		ctx.Write(resp.LoginFail.ToBytes())
		return
	}

	session := sess.Start(ctx)
	session.Set("authenticated", true)

	ctx.Write(resp.SUCCESS.ToBytes())
}

func Logout(ctx iris.Context) {
	session := sess.Start(ctx)

	session.Set("authenticated", false)
	session.Delete("authenticated")
	session.Destroy()

	ctx.Write(resp.SUCCESS.ToBytes())
}

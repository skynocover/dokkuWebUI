package handler

import (
	"os"

	"github.com/kataras/iris/v12"

	"dokkuwebui/src/resp"
	"time"

	"github.com/kataras/iris/v12/sessions"
)

var (
	sessionCookie = "sesson"
	SessionConfig = sessions.Config{Cookie: sessionCookie, Expires: 30 * time.Minute}
	sess          = sessions.New(SessionConfig)
)

func Redirect(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		ctx.Write(resp.SessionExpired.ToBytes())
		return
	}
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

	// Revoke users authentication
	session.Set("authenticated", false)
	// Or to remove the variable:
	session.Delete("authenticated")
	// Or destroy the whole session:
	session.Destroy()

	ctx.Write(resp.SUCCESS.ToBytes())
}

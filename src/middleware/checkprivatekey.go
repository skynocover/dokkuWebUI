package middleware

import (
	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"

	"github.com/kataras/iris/v12"
)

func CheckPrivateKey(ctx iris.Context) {
	if !ssh.Client.CheckSSHKey() {
		ctx.Write(resp.PrivateKeyNotFound.ToBytes())
		return
	}
	ctx.Next()
}

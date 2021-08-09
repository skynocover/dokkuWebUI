package handler

import (
	"bytes"

	"io"

	"github.com/kataras/iris/v12"

	"dokkuwebui/src/resp"
	"dokkuwebui/src/ssh"
)

func Upload(ctx iris.Context) {
	if ssh.Client.CheckSSHKey() {
		ctx.Write(resp.ErrorSSHKeyAlreadyExist.ToBytes())
		return
	}

	f, _, err := ctx.FormFile("uploadfile")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	if err := ssh.Client.Init(buf.Bytes()); err != nil {
		ctx.Write(resp.SystemError.ToBytesWithErr(err))
		return
	}
	ctx.Write(resp.SUCCESS.ToBytes())
}

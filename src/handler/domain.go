package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	if os.Getenv("CF_TOKEN") != "" && os.Getenv("CH_ZONEID") != "" {
		for i := range domains {
			if err := postRecord(domains[i]); err != nil {
				ctx.Write(resp.ErrorCloudFlareFail.ToBytesWithErr(err))
				return
			}
		}
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
	if os.Getenv("CF_TOKEN") != "" && os.Getenv("CH_ZONEID") != "" {
		for i := range domains {
			delRecord(domains[i])
		}
	}

	ctx.Write(resp.SUCCESS.ToBytes())
}

/////////////////////////////// global ////////////////////////////////
func DomainGlobal(ctx iris.Context) {
	r, err := ssh.Client.Run("dokku domains:report --global")
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

/////////////////////// function ////////////////////////////////////

const CFDomain = "https://api.cloudflare.com/client/v4/zones"

func postRecord(name string) (err error) {
	var record struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		Ttl     int    `json:"ttl"`
		Proxied bool   `json:"proxied"`
	}
	record.Type = "A"
	record.Name = strings.Replace(name, os.Getenv("CF_DOMAIN"), "", 1)
	record.Content = os.Getenv("SSH_SERVER")
	record.Ttl = 1
	record.Proxied = true

	jsonbody, _ := json.Marshal(&record)
	sitemap, err := sendToCF("POST", jsonbody, "")
	if err != nil {
		return err
	}

	var result struct {
		Success bool `json:"success"`
		Errors  []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err = json.Unmarshal(sitemap, &result); err != nil {
		return
	}

	if len(result.Errors) > 0 {
		return fmt.Errorf(result.Errors[0].Message)
	}
	return
}

func delRecord(name string) (err error) {
	sitemap, err := sendToCF("GET", []byte{}, "")
	if err != nil {
		return
	}

	var result struct {
		Success bool `json:"success"`
		Errors  []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
		Result []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Content string `json:"content"`
		} `json:"result"`
	}
	if err = json.Unmarshal(sitemap, &result); err != nil {
		return
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf(result.Errors[0].Message)
	}

	for i := range result.Result {
		if result.Result[i].Name == name && result.Result[i].Content == os.Getenv("SSH_SERVER") {
			_, err = sendToCF("DELETE", []byte{}, "/"+result.Result[i].ID)
			if err != nil {
				log.Fatal(err)
				return
			}
			return nil
		}
	}
	return nil
}

func sendToCF(method string, jsonbody []byte, param string) ([]byte, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/%s/dns_records%s", CFDomain, os.Getenv("CH_ZONEID"), param)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonbody))
	if err != nil {
		return nil, err
	}
	token := fmt.Sprintf("Bearer %s", os.Getenv("CF_TOKEN"))
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

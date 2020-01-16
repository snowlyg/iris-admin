package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/kataras/iris/v12"
)

const Key = "2ovULpFZ6l7wzFN"

func Payload(ctx iris.Context) {
	signature := ctx.GetHeader("X-Hub-Signature")
	body, _ := ioutil.ReadAll(ctx.Request().Body)
	//res = []byte(c.GetString(":payload"))
	mac := hmac.New(sha1.New, []byte(Key))
	mac.Write(body)
	s := hex.EncodeToString(mac.Sum(nil))

	if signature == "sha1="+s {
		if err := cmd("cd", "", []string{"/root/go/src/IrisAdminApi"}); err != nil {
			ctx.Application().Logger().Errorf("cmd %v", err)
		}

		if err := cmd("git", "", []string{"pull"}); err != nil {
			ctx.Application().Logger().Errorf("cmd %v", err)
		}

		if err := cmd("go", "", []string{"build"}); err != nil {
			ctx.Application().Logger().Errorf("cmd %v", err)
		}

		if err := cmd("supervisorctl", "", []string{"restart", "iris-admin-api"}); err != nil {
			ctx.Application().Logger().Errorf("cmd %v", err)
		}
	}
}

func cmd(action, input string, arg []string) error {
	cmd := exec.Command(action, arg...)
	if len(input) > 0 {
		cmd.Stdin = strings.NewReader(input)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

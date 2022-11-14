package cmd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "transctl",
	Short: "transctl is a very fast translate command tool",
	Long:  `transctl is a very fast translate command tool`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		raw := fmt.Sprint(strings.Join(args, " "))
		tmp := fmt.Sprintf("%s%s%v%s", id, raw, salt, secret)
		sign := md5.Sum([]byte(tmp))
		param := fmt.Sprintf("q=%s&from=auto&to=zh&appid=%s&salt=%v&sign=%x", raw, id, salt, sign)
		req, err := http.NewRequest("GET", baiduApi+"?"+param, nil)
		if err != nil {
			fmt.Print(req)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Print(err.Error())
			fmt.Print("\n")
		}
		if res.StatusCode != 200 {
			fmt.Printf("request failed, code: %v", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Print(err.Error())
			fmt.Print("\n")
		}
		var newBody responseBody
		if err := json.Unmarshal(body, &newBody); err != nil {
			fmt.Print(err.Error())
		}
		if newBody.ErrorCode != "0" {
			fmt.Print(newBody.ErrorMsg)
		}
		fmt.Printf("%s\n%s", newBody.TransResult[0].Src, newBody.TransResult[0].Dst)
	},
}

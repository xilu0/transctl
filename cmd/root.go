package cmd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

func checkLanguange(Q string) string {
	for _, c := range Q {
		if c > 1000 {
			return "en"
		}
	}
	return "zh"
}

var RootCmd = &cobra.Command{
	Use:   "transctl",
	Short: "transctl is a very fast translate command tool",
	Long:  `transctl is a very fast translate command tool`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		result, err := translate(args)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print(result)
	},
}

func translate(args []string) (string, error) {
	raw := strings.Join(args, " ")
	tmp := fmt.Sprintf("%s%s%v%s", id, raw, salt, secret)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(tmp)))
	params := map[string]string{
		"q":     raw,
		"from":  "auto",
		"to":    checkLanguange(raw),
		"appid": id,
		"salt":  fmt.Sprint(salt),
		"sign":  sign,
	}
	url, err := buildUrlParams(baiduApi, params)
	if err != nil {
		return "", fmt.Errorf("build url params error: %v", err)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("new request error: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do http request error: %v", err)
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("request failed, code: %v", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response body error: %v", err)
	}
	newBody := new(responseBody)
	if err := json.Unmarshal(body, newBody); err != nil {
		return "", fmt.Errorf("unmarshal body error: %v", err)
	}
	if newBody == nil && newBody.ErrorCode != "" {
		return "", fmt.Errorf("translate failed: %v", newBody.ErrorMsg)
	}
	if len(newBody.TransResult) == 0 {
		return "", fmt.Errorf("translate result empty: %v", newBody.ErrorMsg)
	}
	result := newBody.TransResult[0]
	if result.Dst == "" {
		return "", fmt.Errorf("trans failed, dst is empy: %v", newBody.ErrorMsg)
	}
	return "", fmt.Errorf("%s\n%s", result.Src, result.Dst)
}

func buildUrlParams(userUrl string, params map[string]string) (string, error) {
	parsedUrl, err := url.Parse(userUrl)
	if err != nil {
		return "", err
	}
	parsedQuery, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return "", err
	}
	for key, value := range params {
		parsedQuery.Set(key, value)
	}
	return addQueryParams(parsedUrl, parsedQuery), nil
}

func addQueryParams(parsedUrl *url.URL, parsedQuery url.Values) string {
	return strings.Join([]string{strings.Replace(parsedUrl.String(), "?"+parsedUrl.RawQuery, "", -1), parsedQuery.Encode()}, "?")
}

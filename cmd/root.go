package cmd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
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
func GetConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home directory error: %v", err)
	}
	if home == "" {
		if home, err = os.Getwd(); err != nil {
			return "", fmt.Errorf("get user directory error: %v", err)
		}
	}
	var osType = runtime.GOOS
	switch osType {
	case "windows":
		return home + "\\.transctl", nil
	case "linux":
		return home + "/.transctl", nil
	default:
		return "", fmt.Errorf("unknow os type: %v", osType)
	}
}

func GetConfigPath() (string, error) {
	dir, err := GetConfigDirectory()
	if err != nil {
		return "", fmt.Errorf("get config path error: %v", err)
	}
	var osType = runtime.GOOS
	switch osType {
	case "windows":
		return dir + "\\config.json", nil
	case "linux":
		return dir + "/config.json", nil
	default:
		return "", fmt.Errorf("unknow os type: %v", osType)
	}
}

func initConfig() error {
	var accountId string
	var accountSecret string
	fmt.Print("please input your account id: ")
	if _, err := fmt.Scan(&accountId); err != nil {
		return fmt.Errorf("get input error: %v", err)
	}
	fmt.Print("please input your account secret: ")
	if _, err := fmt.Scan(&accountSecret); err != nil {
		return fmt.Errorf("get input error: %v", err)
	}
	var auth = Auth{Id: accountId, Secret: accountSecret}
	authByte, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("marshal auto error: %v", err)
	}
	configDirectory, err := GetConfigDirectory()
	if err != nil {
		return err
	}
	f, err := os.Stat(configDirectory)
	if f != nil && !f.IsDir() {
		if err := os.Remove(configDirectory); err != nil {
			return fmt.Errorf("remove config directory error: %v", err)
		}
	}
	configPath, err := GetConfigDirectory()
	if err != nil {
		return err
	}
	if err != nil {
		if err := os.Mkdir(configDirectory, 0700); err != nil {
			return fmt.Errorf("create config directory error: %v", err)
		}
	}
	if _, err := os.Create(configPath); err != nil {
		return fmt.Errorf("create config error: %v", err)
	}
	if err := os.WriteFile(configPath, authByte, os.FileMode(0600)); err != nil {
		return fmt.Errorf("write config error: %v", err)
	}
	return nil
}

var RootCmd = &cobra.Command{
	Use:   "transctl",
	Short: "transctl is a very fast translate command tool",
	Long:  `transctl is a very fast translate command tool`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		init, err := cmd.Flags().GetBool("init")
		if err != nil {
			fmt.Printf("get init arg error: %v", err)
		}
		if init {
			if err := initConfig(); err != nil {
				fmt.Printf("init config error: %v", err)
			}
			return
		}
		auth, err := getAuth()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		result, err := translate(args, auth)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Print(result)
	},
}

func getAuth() (*Auth, error) {
	file, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read auth config file error: %v", err)
	}
	auth := new(Auth)
	if err := json.Unmarshal(b, auth); err != nil {
		return nil, fmt.Errorf("unmarshal auth error: %v", err)
	}
	return auth, nil
}

func translate(args []string, auth *Auth) (string, error) {
	raw := strings.Join(args, " ")
	tmp := fmt.Sprintf("%s%s%v%s", auth.Id, raw, salt, auth.Secret)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(tmp)))
	params := map[string]string{
		"q":     raw,
		"from":  "auto",
		"to":    checkLanguange(raw),
		"appid": auth.Id,
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

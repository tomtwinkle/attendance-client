package config

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func (c *config) inputTenant() (string, error) {
	validate := func(input string) error {
		if input == "" {
			return errors.New("テナントコードは必須です")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("テナントコードを入力してください"),
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}

func (c *config) inputOBCiD() (string, error) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("OBCiDは数字です")
		}
		if input == "" {
			return errors.New("OBCiDは必須です")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("OBCiDを入力してください"),
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}

func (c *config) inputPassword() (string, error) {
	validate := func(input string) error {
		if input == "" {
			return errors.New("パスワードは必須です")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("パスワードを入力してください"),
		Validate: validate,
		Mask:     '*',
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}

func (c *config) inputUseSlack() (bool, error) {
	validate := func(input string) error {
		if input == "y" || input == "n" {
			return nil
		}
		return errors.New("please input [y or n]")
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Slack Post設定を行いますか？ [y/n]"),
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return false, err
	}
	return result == "y", nil
}

func (c *config) inputSlackToken() (string, error) {
	const slackAppsUrl = "https://api.slack.com/apps"
	validate := func(input string) error {
		if strings.HasPrefix(input, "xoxp-") {
			return nil
		}
		return errors.New("slack tokenを入力してください")
	}
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s",
			"slack tokenを入力してください。",
			"1. [Create New App] を選択",
			"2. [OAuth & Permissions] を選択",
			"3. [User Token Scopes] を選択 [channels:read, chat:write] を追加",
			"4. [Install to Workspace] を選択",
			"5. OAuth Access Token [Copy] を選択",
		),
		Validate: validate,
	}
	if err := browserOpen(slackAppsUrl); err != nil {
		return "", err
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}

func (c *config) inputSlackChannel() (string, error) {
	validate := func(input string) error {
		if input == "" {
			return errors.New("チャンネル名は必須です")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s",
			"投稿したいチャンネル名(#を除いた)を入力してください",
			"チャンネルにメッセージを投稿するためにはチャンネルにアプリを追加する必用があります。",
			"1. [(i) チャンネル詳細] を選択",
			"2. [... その他] を選択",
			"3. [アプリを追加する] を選択",
			"4. 作成したアプリを検索し、[追加]を選択. [表示する]となっている場合は追加は不要です.",
		),
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}

func browserOpen(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
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
	validate := func(input string) error {
		if strings.HasPrefix(input, "xoxp-") {
			return nil
		}
		return errors.New("slack tokenを入力してください")
	}
	prompt := promptui.Prompt{
		Label:    "slack tokenを入力してください",
		Validate: validate,
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
		Label:    "投稿したいチャンネル名を入力してください",
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}

//nolint
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

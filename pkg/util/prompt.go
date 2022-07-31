/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yifan-gu/blueNote/pkg/config"
)

func PromptConfirmation(cfg *config.GlobalConfig, prompt string) (bool, error) {
	if cfg.PromptNoToAll {
		return false, nil
	}
	if cfg.PromptYesToAll {
		return true, nil
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n/yes-to-(a)ll/n(o)ne]: ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "a" || response == "all" {
			cfg.PromptYesToAll = true
			return true, nil
		}
		if response == "o" || response == "none" {
			cfg.PromptNoToAll = true
			return false, nil
		}
		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}
}

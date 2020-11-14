package cmd

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var handlerEntrypointCmd = &cobra.Command{
	Use:   "handler-entrypoint",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dispatchURICommand(cmd, args)
	},
}

func dispatchURICommand(cmd *cobra.Command, args []string) error {
	if runtime.GOOS != "windows" {
		return errors.New("only supported on Windows")
	}

	if len(args) < 1 {
		return errors.New("no arguments given")
	}

	params := strings.Split(args[0], ":")
	// TODO: Handle invalid inputs

	switch params[0] {
	case "download":
		uuid := params[1]
		fmt.Printf("Downloading %s...", uuid)
		err := downloadCmd.RunE(cmd, []string{fmt.Sprintf("--uuid=%s", uuid)})
		if err != nil {
			log.Fatalf("download failed: %v", err)
		}
		return nil
	case "configure":
		if params[2] == "token" {
			token := params[3]
			fmt.Printf("Configuring token=%s...", token)
			err := configureCmd.RunE(cmd, []string{fmt.Sprintf("--token=%s", token)})
			if err != nil {
				log.Fatalf("configuring token failed: %v", err)
			}
		}
		return nil
	}

	return errors.New("invalid command")
}

func init() {
	RootCmd.AddCommand(handlerEntrypointCmd)
}

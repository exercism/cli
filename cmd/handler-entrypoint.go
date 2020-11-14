package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func exit(exitCode int) {
	fmt.Println("\nPress ENTER to exit.")
	_, err := fmt.Scanln()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func dispatchURICommand(cmd *cobra.Command, args []string) error {
	if runtime.GOOS != "windows" {
		fmt.Println("only supported on Windows")
		exit(1)
	}

	if len(args) < 1 {
		fmt.Println("no arguments given")
		exit(1)
	}

	params := strings.Split(args[0], ":")
	// TODO: Handle invalid inputs

	switch params[0] {
	case "download":
		uuid := params[1]
		_cmd := exec.Command("exercism", "download", fmt.Sprintf("--uuid=%s", uuid))
		_cmd.Stderr = os.Stderr
		_cmd.Stdout = os.Stdout
		err := _cmd.Run()
		if err != nil {
			exit(1)
		}
		exit(0)
	case "configure":
		if params[1] == "token" {
			token := params[2]
			//err := configureCmd.RunE(cmd, []string{fmt.Sprintf("--token=%s", token)})
			_cmd := exec.Command("exercism", "configure", fmt.Sprintf("--token=%s", token))
			_cmd.Stderr = os.Stderr
			_cmd.Stdout = os.Stdout
			err := _cmd.Run()
			if err != nil {
				exit(1)
			}
			exit(0)
		}
		exit(1)
	}

	return errors.New("invalid command")
}

func init() {
	RootCmd.AddCommand(handlerEntrypointCmd)
}

package cmd

import (
	"errors"
	"log"
	"runtime"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/registry"
)

var installHandlerCmd = &cobra.Command{
	Use:   "install-handler",
	Short: "Install a custom URI handler.",
	Long:  "Install a custom URI handler.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return registerHandler()
	},
}

func registerHandler() error {
	if runtime.GOOS != "windows" {
		return errors.New("Only supported on Windows")
	}

	// TODO: Probably don't need ALL_ACCESS 0xf003f here
	newk, oldk, err := registry.CreateKey(registry.CLASSES_ROOT, "exercism", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	if oldk {
		log.Println("Handler already registered.")
		return nil
	}

	// Indicate that this key declares a custom URI handler
	newk.SetStringValue("URL Protocol", "")

	shellk, _, err := registry.CreateKey(newk, "shell", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	openk, _, err := registry.CreateKey(shellk, "open", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	commandk, _, err := registry.CreateKey(openk, "command", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	// TODO: Check if exercism is in PATH
	// cmd.exe /k "C:\Users\WDAGUtilityAccount\Desktop\bin\exercism.exe handler-entrypoint %1" << registry
	commandk.SetStringValue("", `iex exercism download --uuid=%1`)

	defer commandk.Close()
	defer openk.Close()
	defer shellk.Close()
	defer newk.Close()

	return nil
}

func init() {
	RootCmd.AddCommand(installHandlerCmd)
}

package handlers

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

// Login interactively stores exercism API configuration.
// Delete when nobody is using 1.6.x anymore.
func Login(ctx *cli.Context) {
	msg := `
	*******************************************************************
	DEPRECATED!

	In the future use the 'exercism configure' command to configure:

		exercism configure --key YOUR_API_KEY

	*******************************************************************

	`
	fmt.Printf(msg)

	dir, err := config.Home()
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Join(dir, config.DirExercises)

	bio := bufio.NewReader(os.Stdin)

	fmt.Print("Your Exercism API key (found at http://exercism.io/account):\n> ")
	key, err := bio.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("What is your exercism exercises path?")
	fmt.Printf("Press Enter to select the default (%s):\n> ", dir)
	dir, err = bio.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// overwrite the context
	gSet := flag.NewFlagSet("global", 0)
	gSet.String("config", ctx.GlobalString("config"), ":nodoc:")
	cSet := flag.NewFlagSet("cmd", 0)
	cSet.String("key", key, ":nodoc:")
	cSet.String("dir", dir, ":nodoc:")
	ctx = cli.NewContext(nil, cSet, gSet)
	// call the Configure() handler
	Configure(ctx)
}

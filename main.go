package main

import (
	"errors"
	"log"

	"github.com/jovandeginste/gr/app"
	"github.com/jovandeginste/gr/common"
	"github.com/spf13/cobra"
)

var (
	f app.Fetcher
	d string
)

func main() {
	cmdRoot := &cobra.Command{
		Use:   "gr",
		Short: "Download and update software packages from Git",
	}

	myInit(cmdRoot)

	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}

func myInit(cmdRoot *cobra.Command) {
	cmdRoot.PersistentFlags().StringVar(&d, "installation-root", "~/.gr", "Location where all software is installed")

	cmdFetch := &cobra.Command{
		Use:   "fetch",
		Short: "Download and update software packages from Git",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a git URL")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			f.Destination = common.NewDestination(d)
			f.ParseURL(args[0])
			if err := f.Fetch(); err != nil {
				panic(err)
			}
		},
	}

	cmdRoot.AddCommand(cmdFetch)
}

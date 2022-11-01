package main

import (
	"errors"
	"fmt"

	"github.com/jovandeginste/gr/pkg/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var a = app.New()

func main() {
	cmdRoot := &cobra.Command{
		Use:   "gr",
		Short: "Download and update software packages from Git",
	}

	myInit(cmdRoot)

	if err := cmdRoot.Execute(); err != nil {
		log().Fatal(err)
	}
}

func myInit(cmdRoot *cobra.Command) {
	var d string

	cmdRoot.PersistentFlags().StringVar(&d, "installation-root", "~/.gr", "Location where all software is installed")
	a.SetDestination(d)

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
			parse(args[0])
		},
	}

	cmdFetch.Flags().BoolVar(&a.Fetcher().Retry, "retry", false, "Whether to remove existing version")

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List currently installed software packages",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}

	cmdRoot.AddCommand(cmdFetch)
	cmdRoot.AddCommand(cmdList)
}

func log() *logrus.Logger {
	return a.Logger
}

func parse(u string) {
	if err := a.Fetcher().ParseURL(u); err != nil {
		log().Fatal(err)
	}

	if err := a.Fetcher().Fetch(); err != nil {
		log().Fatal(err)
	}
}

func list() {
	fmt.Println("Currently installed packages:")

	for _, p := range a.List() {
		fmt.Printf("- %s\n", p)
	}
}

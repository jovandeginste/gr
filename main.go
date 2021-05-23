package main

import (
	"errors"
	"fmt"

	"github.com/jovandeginste/gr/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	d string

	a = app.New()
	f = a.Fetcher()
)

func main() {
	cmdRoot := &cobra.Command{
		Use:   "gr",
		Short: "Download and update software packages from Git",
	}

	myInit(cmdRoot)

	a.SetDestination(d)

	if err := cmdRoot.Execute(); err != nil {
		log().Fatal(err)
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
			parse(args[0])
		},
	}
	cmdFetch.Flags().BoolVar(&f.Retry, "retry", false, "Whether to remove existing version")

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
	if err := f.ParseURL(u); err != nil {
		log().Fatal(err)
	}

	if err := f.Fetch(); err != nil {
		log().Fatal(err)
	}
}

func list() {
	fmt.Println("Currently installed packages:")

	for _, p := range a.List() {
		fmt.Printf("- %s\n", p)
	}
}

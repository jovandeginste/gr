package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *cli) listCmd() *cobra.Command {
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List currently installed software packages",
		Run: func(cmd *cobra.Command, args []string) {
			c.list()
		},
	}

	return cmdList
}

func (c *cli) list() {
	fmt.Println("Currently installed packages:")

	for _, p := range c.app.List() {
		fmt.Printf("- %s\n", p)
	}
}

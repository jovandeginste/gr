package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed assets/function.bash
var bashCode string

func (c *cli) envCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Generates code to use this tool with your shell",
		Run: func(cmd *cobra.Command, args []string) {
			c.emitEnv()
		},
	}
	cmd.Flags().String("shell", "auto", "Which shell to emit code for")
	c.viper.BindPFlag("shell", cmd.Flags().Lookup("shell")) //nolint:errcheck

	return cmd
}

func (c *cli) emitEnv() {
	s := c.config.Shell
	if s == "auto" {
		s = os.Getenv("SHELL")
		s = filepath.Base(s)
	}

	switch s {
	case "bash":
		c.emitBashEnv()
	default:
		c.log().Errorf("Unsupported shell '%s'", s)
	}
}

func (c *cli) emitBashEnv() {
	fmt.Println("# Generated with:", strings.Join(os.Args, " "))

	bashCode = strings.ReplaceAll(bashCode, "%ROOT_DIRECTORY%", c.app.Destination.Root)

	fmt.Println(strings.TrimSpace(bashCode))
}

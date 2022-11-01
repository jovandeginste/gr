package main

import (
	"path"

	"github.com/adrg/xdg"
	"github.com/jovandeginste/gr/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cli struct {
	app    *app.App
	cmd    *cobra.Command
	viper  *viper.Viper
	config app.Configuration
}

func (c *cli) myInit() {
	cobra.OnInitialize(c.updateConfig)

	c.cmd = c.rootCmd()
	c.cmd.AddCommand(c.envCmd())
	c.cmd.AddCommand(c.fetchCmd())
	c.cmd.AddCommand(c.listCmd())
}

func (c *cli) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gr",
		Short: "Download and update software packages from Git",
	}
	cmd.PersistentFlags().String("installation-root", "~/.gr", "Location where all software is installed")
	c.viper.BindPFlag("root_directory", cmd.PersistentFlags().Lookup("installation-root")) //nolint:errcheck

	return cmd
}

func (c *cli) initConfig() error {
	c.viper = viper.New()

	c.viper.AddConfigPath(".")
	c.viper.AddConfigPath(path.Join(xdg.ConfigHome, "gr"))
	c.viper.SetConfigName("config") // Register config file name (no extension)

	if err := c.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

func (c *cli) updateConfig() {
	if err := c.viper.Unmarshal(&c.config); err != nil {
		panic(err)
	}

	c.app = app.New(c.config)
}

package main

import "github.com/spf13/cobra"

func (c *cli) fetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch [url]",
		Short: "Download and update software packages from Git",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.fetch(args[0]); err != nil {
				c.log().Error(err)
			}
		},
	}

	cmd.Flags().Bool("retry", false, "Whether to remove existing version")
	c.viper.BindPFlag("retry", cmd.Flags().Lookup("retry")) //nolint:errcheck

	return cmd
}

func (c *cli) fetch(u string) error {
	f, err := c.app.NewFetcher(u)
	if err != nil {
		return err
	}

	return f.Fetch()
}

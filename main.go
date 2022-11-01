package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	c := cli{}

	if err := c.initConfig(); err != nil {
		panic(err)
	}

	c.myInit()

	if err := c.cmd.Execute(); err != nil {
		c.log().Error(err)
	}
}

func (c *cli) log() *logrus.Logger {
	return c.app.Logger
}

package app

import (
	"github.com/jovandeginste/gr/common"
	"github.com/sirupsen/logrus"
)

type App struct {
	Destination *common.Destination
	Preferences *common.Preferences
	Logger      *logrus.Logger
}

func New() *App {
	a := App{
		Logger: logrus.New(),
	}

	return &a
}

func (a *App) Fetcher() *Fetcher {
	f := Fetcher{
		app: a,
	}

	return &f
}

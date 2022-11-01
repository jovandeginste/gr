package app

import (
	"io/ioutil"

	"github.com/jovandeginste/gr/pkg/common"
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

func (a *App) SetDestination(root string) {
	a.Destination = common.NewDestination(root)
}

func (a *App) Fetcher() *Fetcher {
	f := Fetcher{
		app: a,
	}

	return &f
}

func (a *App) List() []string {
	files, err := ioutil.ReadDir(a.Destination.PackagesDir)
	if err != nil {
		return nil
	}

	res := []string{}

	for _, f := range files {
		res = append(res, f.Name())
	}

	return res
}

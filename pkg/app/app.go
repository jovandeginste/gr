package app

import (
	"io/ioutil"

	"github.com/jovandeginste/gr/pkg/common"
	"github.com/sirupsen/logrus"
)

type App struct {
	Destination   *common.Destination
	Preferences   *common.Preferences
	Logger        *logrus.Logger
	Configuration Configuration
}

func New(cfg Configuration) *App {
	a := App{
		Configuration: cfg,
		Logger:        logrus.New(),
	}

	a.SetDestination(cfg.RootDirectory)

	return &a
}

func (a *App) SetDestination(root string) {
	a.Destination = common.NewDestination(root)
}

func (a *App) NewFetcher(url string) (*Fetcher, error) {
	f := Fetcher{
		app: a,
	}
	if err := f.ParseURL(url); err != nil {
		return nil, err
	}

	return &f, nil
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

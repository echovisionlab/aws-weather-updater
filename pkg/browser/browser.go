package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func New() (*rod.Browser, *launcher.Launcher, error) {
	l := launcher.New()
	url, err := l.Launch()
	if err != nil {
		return nil, nil, err
	}
	b := rod.New().ControlURL(url)
	if err = b.Connect(); err != nil {
		return nil, nil, err
	}
	return b, l, nil
}

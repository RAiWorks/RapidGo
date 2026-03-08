package main

import (
	"github.com/RAiWorks/RapidGo/app/providers"
	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/cli"
	"github.com/RAiWorks/RapidGo/core/service"
)

func main() {
	cli.SetBootstrap(func(a *app.App, mode service.Mode) {
		a.Register(&providers.ConfigProvider{})
		a.Register(&providers.LoggerProvider{})
		if mode.Has(service.ModeWeb) || mode.Has(service.ModeAPI) || mode.Has(service.ModeWS) {
			a.Register(&providers.DatabaseProvider{})
		}
		a.Register(&providers.RedisProvider{})
		a.Register(&providers.QueueProvider{})
		if mode.Has(service.ModeWeb) {
			a.Register(&providers.SessionProvider{})
		}
		a.Register(&providers.MiddlewareProvider{Mode: mode})
		a.Register(&providers.RouterProvider{Mode: mode})
	})

	cli.Execute()
}

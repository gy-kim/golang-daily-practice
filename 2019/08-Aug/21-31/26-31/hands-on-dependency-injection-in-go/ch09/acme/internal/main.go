package main

import (
	"context"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/rest"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/config"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/exchange"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/get"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/list"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/register"
)

func main() {
	// bind stop channel to context
	ctx := context.Background()

	// build the exchanger
	exchanger := exchange.NewConverter(config.App)

	// build model layer
	getModel := get.NewGetter(config.App)
	listModel := list.NewLister(config.App)
	registerModel := register.NewRegisterer(config.App, exchanger)

	// start REST server
	server := rest.New(config.App, getModel, listModel, registerModel)

	config.App.Logger().Debug("Starting server.")
	server.Listen(ctx.Done())
}

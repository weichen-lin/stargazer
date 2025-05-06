package main

import (
	"github.com/weichen-lin/stargazer/background"
	"github.com/weichen-lin/stargazer/controller"
	"github.com/weichen-lin/stargazer/infrastructure"
	"github.com/weichen-lin/stargazer/service"
)

func initStarGazer(config *infrastructure.Config) (*controller.Controller, *background.Background) {
	service := service.NewService(config)

	controller := controller.NewController(service)
	background := background.NewBackground(service, config)

	return controller, background
}

package internal

import "time4book/internal/app"

type App struct {
	App *app.Facade
}

func New() *App {
	appFacade := app.New()

	return &App{
		App: appFacade,
	}
}

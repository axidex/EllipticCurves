package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/axidex/elliptic/internal/gui"
	"github.com/axidex/elliptic/internal/logger"
)

func main() {
	// Создаем новое приложение
	a := app.New()

	appLogger := logger.NewGUILogger()

	guiApp := gui.NewAppGui(a, appLogger)

	guiApp.Run()
}

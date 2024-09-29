package gui

import "fyne.io/fyne/v2"

// Функция для обновления размера окна
func (app *AppGui) updateWindowSize() {
	app.w.Resize(fyne.NewSize(app.width, app.height))
}

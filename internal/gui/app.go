package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/axidex/Unknown/pkg/logger"
)

type App interface {
	Run()
}

type AppGui struct {
	fyneApp                              fyne.App
	logger                               logger.Logger
	privateKeyEntry, publicKeyEntry      *widget.Entry
	openText, closedText                 *widget.Entry
	encrypt, decrypt, generateKeysButton *widget.Button
	w                                    fyne.Window

	width, height float32
}

func NewAppGui(app fyne.App, logger logger.Logger) AppGui {
	appGui := AppGui{
		fyneApp: app,
		logger:  logger,
	}
	appGui.w = app.NewWindow("ECIES app")

	appGui.initEntry()
	appGui.initButtons()

	appGui.width = float32(800)
	appGui.height = float32(1000)

	return appGui
}

func (app *AppGui) initEntry() {
	// Создаем поля ввода
	app.privateKeyEntry = widget.NewMultiLineEntry()
	app.privateKeyEntry.SetPlaceHolder("Private Key")

	app.publicKeyEntry = widget.NewMultiLineEntry()
	app.publicKeyEntry.SetPlaceHolder("Public Key")

	app.openText = widget.NewMultiLineEntry()
	app.openText.SetPlaceHolder("Open Text")

	app.closedText = widget.NewMultiLineEntry()
	app.closedText.SetPlaceHolder("Closed Text")

}

func (app *AppGui) initButtons() {
	// Создаем кнопки
	app.encrypt = widget.NewButton("Encrypt", app.encryptData)
	app.decrypt = widget.NewButton("Decrypt", app.decryptData)

	app.generateKeysButton = widget.NewButton("Generate Keys", app.generateKeys)
}

func (app *AppGui) Run() {

	app.height = app.privateKeyEntry.MinSize().Height +
		app.publicKeyEntry.MinSize().Height +
		app.openText.MinSize().Height +
		app.closedText.MinSize().Height +
		app.encrypt.MinSize().Height +
		app.decrypt.MinSize().Height +
		app.generateKeysButton.MinSize().Height + 50

	// Обработчики для обновления размера окна при изменении текста
	app.privateKeyEntry.OnChanged = func(s string) { app.updateWindowSize() }
	app.publicKeyEntry.OnChanged = func(s string) { app.updateWindowSize() }
	app.openText.OnChanged = func(s string) { app.updateWindowSize() }
	app.closedText.OnChanged = func(s string) { app.updateWindowSize() }

	// Используем контейнер с вертикальной компоновкой для размещения элементов
	content := container.NewVBox(
		app.privateKeyEntry,
		app.publicKeyEntry,
		app.openText,
		app.closedText,
		app.encrypt,
		app.decrypt,
		app.generateKeysButton,
	)

	// Устанавливаем содержимое окна
	app.w.SetContent(content)

	// Задаем минимальный размер окна
	app.w.Resize(fyne.NewSize(app.width, app.height))

	// Показываем окно
	app.w.ShowAndRun()
}

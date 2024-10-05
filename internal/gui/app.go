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
	curveInfoEntry                       *widget.Entry
	encrypt, decrypt, generateKeysButton *widget.Button
	w                                    fyne.Window

	width, height float32
}

func NewAppGui(app fyne.App, logger logger.Logger) AppGui {
	appGui := AppGui{
		fyneApp: app,
		logger:  logger,
	}
	app.Settings().SetTheme(EllipticTheme{})
	appGui.w = app.NewWindow("ECIES")

	appGui.initEntry()
	appGui.initButtons()

	appGui.width = float32(680)
	//appGui.height = float32(1000)

	return appGui
}

func (app *AppGui) initEntry() {

	app.privateKeyEntry = widget.NewMultiLineEntry()
	initEntry(app.privateKeyEntry, "Private Key", 8)

	app.publicKeyEntry = widget.NewMultiLineEntry()
	initEntry(app.publicKeyEntry, "Public Key", 7)

	app.openText = widget.NewMultiLineEntry()
	initEntry(app.openText, "Open Text", 3)

	app.closedText = widget.NewMultiLineEntry()
	initEntry(app.closedText, "Closed Text", 3)

	app.curveInfoEntry = widget.NewMultiLineEntry()
	initEntry(app.curveInfoEntry, "Curve Info", 6)
	app.curveInfoEntry.Disable()

}

func initEntry(entry *widget.Entry, name string, numberOfLines int) {
	entry.SetPlaceHolder(name)
	entry.SetMinRowsVisible(numberOfLines)
	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = container.ScrollHorizontalOnly
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

	content := container.NewVBox(
		app.privateKeyEntry,
		app.publicKeyEntry,
		app.openText,
		app.closedText,
		app.curveInfoEntry,
		app.encrypt,
		app.decrypt,
		app.generateKeysButton,
	)

	app.w.SetContent(content)

	size := fyne.NewSize(app.width, app.height)
	app.w.Resize(size)
	app.w.SetFixedSize(true)

	// Показываем окно
	app.w.ShowAndRun()
}

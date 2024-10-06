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
	curveInfoEntry, eciesInfo            *widget.Entry
	encrypt, decrypt, generateKeysButton *widget.Button
	selectCurve, selectParams            *widget.Select
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
	appGui.initSelect()

	appGui.width = float32(1200)
	//appGui.height = float32(1000)

	return appGui
}

func (app *AppGui) initSelect() {
	app.selectParams = widget.NewSelect(ParamNames, func(string) {})
	app.selectParams.SetSelectedIndex(0)

	app.selectCurve = widget.NewSelect(CurveNames, func(string) {})
	app.selectCurve.SetSelectedIndex(0)
}

func (app *AppGui) initEntry() {

	app.privateKeyEntry = widget.NewMultiLineEntry()
	initEntry(app.privateKeyEntry, "Private Key", 8)

	app.publicKeyEntry = widget.NewMultiLineEntry()
	initEntry(app.publicKeyEntry, "Public Key", 8)

	app.openText = widget.NewMultiLineEntry()
	initEntry(app.openText, "Open Text", 3)

	app.closedText = widget.NewMultiLineEntry()
	initEntry(app.closedText, "Closed Text", 3)

	app.curveInfoEntry = widget.NewMultiLineEntry()
	initEntry(app.curveInfoEntry, "Curve Info", 7)
	app.curveInfoEntry.Disable()

	app.eciesInfo = widget.NewMultiLineEntry()
	initEntry(app.eciesInfo, "ECIES Info", 7)
	app.eciesInfo.Disable()

}

func initEntry(entry *widget.Entry, name string, numberOfLines int) {
	entry.SetPlaceHolder(name)
	entry.SetMinRowsVisible(numberOfLines)
	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = container.ScrollHorizontalOnly
}

func (app *AppGui) initButtons() {
	app.encrypt = widget.NewButton("Encrypt", app.encryptData)
	app.decrypt = widget.NewButton("Decrypt", app.decryptData)

	app.generateKeysButton = widget.NewButton("Generate Keys", app.generateKeys)
}

func (app *AppGui) Run() {

	app.height = app.privateKeyEntry.MinSize().Height +
		app.closedText.MinSize().Height +
		app.eciesInfo.MinSize().Height +
		app.decrypt.MinSize().Height + 40

	leftContainer := container.NewVBox(
		app.privateKeyEntry,
		app.closedText,
		app.eciesInfo,
		app.decrypt,
		app.selectParams,
	)

	rightContainer := container.NewVBox(
		app.publicKeyEntry,
		app.openText,
		app.curveInfoEntry,
		app.encrypt,
		app.selectCurve,
	)

	topContainer := container.NewGridWithColumns(2, leftContainer, rightContainer)

	bottomContainer := container.NewGridWrap(
		fyne.NewSize(app.width, app.generateKeysButton.MinSize().Height),
		app.generateKeysButton,
	)

	appContainer := container.NewVBox(topContainer, bottomContainer)

	app.w.SetContent(appContainer)

	size := fyne.NewSize(app.width, app.height)
	app.w.Resize(size)
	app.w.SetFixedSize(true)

	// Показываем окно
	app.w.ShowAndRun()
}

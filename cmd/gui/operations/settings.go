package operations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/internal/i18n"
	"github.com/lcpu-club/kube-dify-gui/internal/ui"
)

// ShowSettings displays the settings screen
func ShowSettings(window fyne.Window, translator *i18n.Translator, onBack func(), onLanguageChanged func(string), onThemeChanged func(bool)) {
	// Create UI components
	components := ui.New(translator)

	// Create language selector
	languageLabel := widget.NewLabel(translator.Translate("settings.language"))
	languageSelector := components.CreateLanguageSelector(func(lang string) {
		if onLanguageChanged != nil {
			onLanguageChanged(lang)
		}
	})

	// Create theme selector
	themeLabel := widget.NewLabel(translator.Translate("settings.theme"))
	themeOptions := []string{
		translator.Translate("settings.theme.light"),
		translator.Translate("settings.theme.dark"),
	}

	// Determine current theme
	isDark := fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantDark
	initialTheme := 0
	if isDark {
		initialTheme = 1
	}

	themeSelector := widget.NewSelect(themeOptions, func(selected string) {
		if selected == translator.Translate("settings.theme.dark") {
			if onThemeChanged != nil {
				onThemeChanged(true)
			}
		} else {
			if onThemeChanged != nil {
				onThemeChanged(false)
			}
		}
	})
	themeSelector.SetSelectedIndex(initialTheme)

	// Create back button
	backButton := components.CreateBackButton(onBack)

	// Create layout
	settingsForm := container.NewGridWithColumns(2,
		languageLabel, languageSelector,
		themeLabel, themeSelector,
	)

	// Create header
	header := components.CreateHeader(translator.Translate("settings.title"))

	// Create main content area
	mainContent := container.NewVBox(
		settingsForm,
	)

	// Create card for main content
	contentCard := components.CreateCard("", mainContent)

	// Create layout
	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		contentCard,
		backButton,
	)

	window.SetContent(container.NewPadded(content))
}

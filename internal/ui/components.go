package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/internal/i18n"
)

// Components provides reusable UI components
type Components struct {
	translator *i18n.Translator
}

// New creates a new UI components provider
func New(translator *i18n.Translator) *Components {
	return &Components{
		translator: translator,
	}
}

// Translate is a shorthand for translator.Translate
func (c *Components) Translate(id string) string {
	return c.translator.Translate(id)
}

// CreateCard creates a card container with a title and content
func (c *Components) CreateCard(title string, content fyne.CanvasObject) *fyne.Container {
	if title == "" {
		// If no title, just return the content without a separator
		return container.NewVBox(content)
	}

	// If there's a title, add a title label and separator
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	return container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		content,
	)
}

// CreateInfoCard creates an information card with an icon and text
func (c *Components) CreateInfoCard(icon fyne.Resource, text string) *fyne.Container {
	iconObj := widget.NewIcon(icon)

	// Create a text label that will fill available space
	textLabel := widget.NewLabel(text)
	textLabel.Wrapping = fyne.TextWrapWord

	// Use a border layout to position the icon on the left and let the text fill the rest
	return container.NewBorder(
		nil, nil, iconObj, nil,
		container.NewPadded(textLabel), // Padding around the text for better appearance
	)
}

// CreateOperationButton creates a button for an operation with an icon, label, and tooltip
func (c *Components) CreateOperationButton(icon fyne.Resource, labelID string, tooltipID string, onTapped func()) *widget.Button {
	button := widget.NewButtonWithIcon(c.Translate(labelID), icon, onTapped)
	button.Importance = widget.MediumImportance

	// Set tooltip
	tooltip := c.Translate(tooltipID)
	if tooltip != tooltipID { // Only set if translation exists
		button.Importance = widget.MediumImportance
		// In Fyne v2, we can't directly set tooltips on buttons
		// We'll handle this in the UI layout by wrapping buttons in containers with tooltips
	}

	return button
}

// CreateBackButton creates a standard back button
func (c *Components) CreateBackButton(onTapped func()) *widget.Button {
	return widget.NewButtonWithIcon(c.Translate("common.back"), theme.NavigateBackIcon(), onTapped)
}

// CreateConfirmationDialog creates a confirmation dialog
func (c *Components) CreateConfirmationDialog(title, message string, onConfirm func(), parent fyne.Window) *widget.PopUp {
	// Create dialog content
	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord

	// Create buttons
	confirmButton := widget.NewButtonWithIcon(c.Translate("common.confirm"), theme.ConfirmIcon(), func() {
		if onConfirm != nil {
			onConfirm()
		}
	})
	confirmButton.Importance = widget.HighImportance

	cancelButton := widget.NewButtonWithIcon(c.Translate("common.cancel"), theme.CancelIcon(), nil)

	// Create dialog
	content := container.NewVBox(
		widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		messageLabel,
		container.NewHBox(
			cancelButton,
			confirmButton,
		),
	)

	popup := widget.NewModalPopUp(content, parent.Canvas())

	// Set button actions
	cancelButton.OnTapped = func() {
		popup.Hide()
	}
	confirmButton.OnTapped = func() {
		popup.Hide()
		if onConfirm != nil {
			onConfirm()
		}
	}

	return popup
}

// CreateStatusBar creates a status bar with a label
func (c *Components) CreateStatusBar(initialStatus string) (*fyne.Container, *widget.Label) {
	statusLabel := widget.NewLabel(initialStatus)

	return container.NewHBox(
		statusLabel,
	), statusLabel
}

// CreateSettingsButton creates a settings button
func (c *Components) CreateSettingsButton(onTapped func()) *widget.Button {
	return widget.NewButtonWithIcon("", theme.SettingsIcon(), onTapped)
}

// CreateLogoutButton creates a logout button
func (c *Components) CreateLogoutButton(onTapped func()) *widget.Button {
	return widget.NewButtonWithIcon(c.Translate("dashboard.logout"), theme.LogoutIcon(), onTapped)
}

// CreateHelpIcon creates a help icon with a tooltip
func (c *Components) CreateHelpIcon(helpText string) *fyne.Container {
	icon := widget.NewIcon(theme.InfoIcon())
	iconContainer := container.NewCenter(icon)
	// In Fyne v2, we can't directly set tooltips on containers
	// We'll handle this in the UI layout

	return iconContainer
}

// CreateTooltip wraps a widget with a tooltip using a custom tooltip widget
func (c *Components) CreateTooltip(content fyne.CanvasObject, tooltip string) *fyne.Container {
	return container.NewStack(content)
}

// CreateHeader creates a header with a title and optional buttons
func (c *Components) CreateHeader(title string, buttons ...*widget.Button) *fyne.Container {
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	rightContent := container.NewHBox()
	for _, button := range buttons {
		rightContent.Add(button)
	}

	return container.NewBorder(
		nil, nil, nil, rightContent,
		titleLabel,
	)
}

// CreateLanguageSelector creates a language selector
func (c *Components) CreateLanguageSelector(onChanged func(string)) *widget.Select {
	languages := c.translator.GetAvailableLanguages()
	languageNames := make([]string, len(languages))

	for i, lang := range languages {
		languageNames[i] = c.translator.GetLanguageName(lang)
	}

	selector := widget.NewSelect(languageNames, func(selected string) {
		// Find the language code for the selected name
		var langCode string
		for _, lang := range languages {
			if c.translator.GetLanguageName(lang) == selected {
				langCode = lang
				break
			}
		}

		if langCode != "" && onChanged != nil {
			onChanged(langCode)
		}
	})

	// Set the current language
	currentLang := c.translator.GetLanguage()
	selector.SetSelected(c.translator.GetLanguageName(currentLang))

	return selector
}

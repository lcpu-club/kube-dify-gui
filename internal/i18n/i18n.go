package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed translations/*.json
var translationFS embed.FS

// Translator handles translations for the application
type Translator struct {
	bundle       *i18n.Bundle
	localizer    *i18n.Localizer
	currentLang  string
	app          fyne.App
	translations map[string]map[string]string
}

// New creates a new translator
func New(a fyne.App) (*Translator, error) {
	if a == nil {
		a = app.New()
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load translations from embedded files
	entries, err := translationFS.ReadDir("translations")
	if err != nil {
		// If we can't read the directory, create an empty translator
		// This allows the app to run without translations
		return &Translator{
			bundle:       bundle,
			localizer:    i18n.NewLocalizer(bundle, "en"),
			currentLang:  "en",
			app:          a,
			translations: make(map[string]map[string]string),
		}, nil
	}

	translations := make(map[string]map[string]string)

	// Load each translation file
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			data, err := translationFS.ReadFile(fmt.Sprintf("translations/%s", entry.Name()))
			if err != nil {
				continue
			}

			// Extract language code from filename (e.g., "en.json" -> "en")
			lang := strings.TrimSuffix(entry.Name(), ".json")

			// Parse the JSON data
			var langMap map[string]string
			if err := json.Unmarshal(data, &langMap); err != nil {
				continue
			}

			translations[lang] = langMap

			// Add to bundle
			_, err = bundle.ParseMessageFileBytes(data, entry.Name())
			if err != nil {
				continue
			}
		}
	}

	// Get preferred language from app preferences, default to English
	lang := a.Preferences().StringWithFallback("language", "en")

	return &Translator{
		bundle:       bundle,
		localizer:    i18n.NewLocalizer(bundle, lang),
		currentLang:  lang,
		app:          a,
		translations: translations,
	}, nil
}

// SetLanguage sets the current language
func (t *Translator) SetLanguage(lang string) {
	t.currentLang = lang
	t.localizer = i18n.NewLocalizer(t.bundle, lang)
	t.app.Preferences().SetString("language", lang)
}

// GetLanguage returns the current language
func (t *Translator) GetLanguage() string {
	return t.currentLang
}

// Translate translates a message ID to the current language
func (t *Translator) Translate(id string) string {
	// Try to get from the localizer
	msg, err := t.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})

	if err == nil && msg != "" {
		return msg
	}

	// Fall back to direct lookup in translations map
	if langMap, ok := t.translations[t.currentLang]; ok {
		if translation, ok := langMap[id]; ok {
			return translation
		}
	}

	// Return the ID as a last resort
	return id
}

// GetAvailableLanguages returns the list of available languages
func (t *Translator) GetAvailableLanguages() []string {
	languages := make([]string, 0, len(t.translations))
	for lang := range t.translations {
		languages = append(languages, lang)
	}
	return languages
}

// GetLanguageName returns the display name for a language code
func (t *Translator) GetLanguageName(code string) string {
	switch code {
	case "en":
		return "English"
	case "zh":
		return "中文"
	default:
		return code
	}
}

package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// Settings represents the application settings
type Settings struct {
	app fyne.App
}

// New creates a new settings manager
func New(a fyne.App) *Settings {
	if a == nil {
		a = app.New()
	}
	return &Settings{
		app: a,
	}
}

// SaveToken saves the user token to preferences
func (s *Settings) SaveToken(token string) {
	s.app.Preferences().SetString("user_token", token)
}

// LoadToken loads the user token from preferences
func (s *Settings) LoadToken() string {
	return s.app.Preferences().StringWithFallback("user_token", "")
}

// ClearToken clears the user token from preferences
func (s *Settings) ClearToken() {
	s.app.Preferences().RemoveValue("user_token")
}

// SetRememberLogin sets whether to remember login
func (s *Settings) SetRememberLogin(remember bool) {
	s.app.Preferences().SetBool("remember_login", remember)
}

// GetRememberLogin gets whether to remember login
func (s *Settings) GetRememberLogin() bool {
	return s.app.Preferences().BoolWithFallback("remember_login", false)
}

// SetLanguage sets the preferred language
func (s *Settings) SetLanguage(lang string) {
	s.app.Preferences().SetString("language", lang)
}

// GetLanguage gets the preferred language
func (s *Settings) GetLanguage() string {
	return s.app.Preferences().StringWithFallback("language", "en")
}

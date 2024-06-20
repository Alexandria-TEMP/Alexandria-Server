package forms

type TokenRefreshForm struct {
	RefreshToken string `json:"refreshToken"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *TokenRefreshForm) IsValid() bool {
	return form.RefreshToken != ""
}

package nbpapi

// SetLang function
func SetLang(lang string) {
	if lang == "pl" {
		l = langTexts["pl"]
	} else {
		l = langTexts["en"]
	}
}

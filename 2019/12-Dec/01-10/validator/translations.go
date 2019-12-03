package validator

import ut "github.com/go-playground/universal-translator"

type TranslationFunc func(ut ut.Translator, fe FieldError) string

type RegisterTranslationsFunc func(ut ut.Translator) error

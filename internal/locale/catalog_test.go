// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package locale // import "miniflux.app/v2/internal/locale"

import (
	"testing"
)

func TestParserWithInvalidData(t *testing.T) {
	_, err := parseTranslationMessages([]byte(`{`))
	if err == nil {
		t.Fatal(`An error should be returned when parsing invalid data`)
	}
}

func TestParser(t *testing.T) {
	translations, err := parseTranslationMessages([]byte(`{"k": "v"}`))
	if err != nil {
		t.Fatalf(`Unexpected parsing error: %v`, err)
	}

	value, found := translations.singulars["k"]
	if !found {
		t.Fatalf(`The translation %v should contains the defined key`, translations.singulars)
	}

	if value != "v" {
		t.Fatal(`The translation key should contains the defined value`)
	}
}

func TestLoadCatalog(t *testing.T) {
	for language := range AvailableLanguages {
		_, err := loadTranslationFile(language)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAllKeysHaveValue(t *testing.T) {
	for language := range AvailableLanguages {
		messages, err := loadTranslationFile(language)
		if err != nil {
			t.Fatalf(`Unable to load translation messages for language %q`, language)
		}

		if len(messages.singulars) == 0 {
			t.Fatalf(`The language %q doesn't have any messages for singulars`, language)
		}

		if len(messages.plurals) == 0 {
			t.Fatalf(`The language %q doesn't have any messages for plurals`, language)
		}

		for k, v := range messages.singulars {
			if len(v) == 0 {
				t.Errorf(`The key %q for singulars for the language %q has an empty list as value`, k, language)
			}
		}
		for k, v := range messages.plurals {
			if len(v) == 0 {
				t.Errorf(`The key %q for plurals for the language %q has an empty list as value`, k, language)
			}
		}
	}
}

func TestMissingTranslations(t *testing.T) {
	refLang := "en_US"
	references, err := loadTranslationFile(refLang)
	if err != nil {
		t.Fatal(`Unable to parse reference language`)
	}

	for language := range AvailableLanguages {
		if language == refLang {
			continue
		}

		messages, err := loadTranslationFile(language)
		if err != nil {
			t.Fatalf(`Parsing error for language %q`, language)
		}

		for key := range references.singulars {
			if _, found := messages.singulars[key]; !found {
				t.Errorf(`Translation key %q not found in language %q singulars`, key, language)
			}
		}
		for key := range references.plurals {
			if _, found := messages.plurals[key]; !found {
				t.Errorf(`Translation key %q not found in language %q plurals`, key, language)
			}
		}
	}
}

func TestTranslationFilePluralForms(t *testing.T) {
	var numberOfPluralFormsPerLanguage = map[string]int{
		"de_DE":            2,
		"el_EL":            2,
		"en_US":            2,
		"es_ES":            2,
		"fi_FI":            2,
		"fr_FR":            2,
		"hi_IN":            2,
		"id_ID":            1,
		"it_IT":            2,
		"ja_JP":            1,
		"nan_Latn_pehoeji": 1,
		"nl_NL":            2,
		"pl_PL":            3,
		"pt_BR":            2,
		"ro_RO":            3,
		"ru_RU":            3,
		"tr_TR":            2,
		"uk_UA":            3,
		"zh_CN":            1,
		"zh_TW":            1,
	}
	for language := range AvailableLanguages {
		messages, err := loadTranslationFile(language)
		if err != nil {
			t.Fatalf(`Unable to load translation messages for language %q`, language)
		}

		for k, v := range messages.plurals {
			if len(v) != numberOfPluralFormsPerLanguage[language] {
				t.Errorf(`The key %q for the language %q does not have the expected number of plurals, got %d instead of %d`, k, language, len(v), numberOfPluralFormsPerLanguage[language])
			}
		}
	}
}

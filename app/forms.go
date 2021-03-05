package app

import (
	"errors"
	"fmt"
	"github.com/cfc-servers/cfc_suggestions/forms"
)

var Forms = []forms.Form{
	{
		Name: "suggestion",

		Validators: []forms.Validator{
			forms.FieldValidator{}.Accept(
				forms.Field("why").Min(18).Max(1024),
				forms.Field("whyNot").Min(18).Max(1024),
				forms.Field("description").Min(18).Max(1024).Optional(),
				forms.Field("title").Min(6).Max(124),
				forms.Field("image").Max(124).Optional(),
				forms.Field("realm").MustBe("cfc3", "cfcrp", "cfcmc", "cfcrvr", "discord", "other"),
				forms.Field("link").Max(124).Optional(),
				forms.Field("anonymous").Optional(),
			),
		},

		Destinations: []forms.Destination{
			SuggestionsDestination,
		},

		Formatter: forms.DefaultFormatter{Color: 0x34eb5b}.Fields(
			"why", "whyNot", "link",
		),
	},
}

var ErrMissingForm = errors.New("a form with that name did not exist")

func GetForm(name string) (forms.Form, error) {
	for _, form := range Forms {
		if form.Name == name {
			return form, nil
		}
	}

	return forms.Form{}, fmt.Errorf("%w: %v", ErrMissingForm, name)
}

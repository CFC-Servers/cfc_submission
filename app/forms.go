package app

import (
	"errors"
	"fmt"
	"github.com/cfc-servers/cfc_suggestions/forms"
)



var Forms = []forms.Form{
	{
		Name:         "suggestion",
		Validators: []forms.Validator{
			forms.FieldValidator{}.Accept(
				forms.Field("whyNot").Min(10).Max(300),
				forms.Field("why").Min(10).Max(300),
				forms.Field("description").Min(10).Max(300).Optional(),
				forms.Field("image").Max(100).Optional(),
				forms.Field("realm").MustBe(	"cfc3", "cfcrp", "cfcmc", "cfcrvr", "discord", "other"),
			),
		},
		Destinations: []forms.Destination{
			SuggestionsDestination,
		},
	},

}

var ErrMissingForm = errors.New("a form with that name did not exist")

func GetForm(name string) (forms.Form,  error) {
	for _, form := range Forms {
		if form.Name == name {
			return form, nil
		}
	}

	return forms.Form{}, fmt.Errorf("%w: %v", ErrMissingForm, name)
}

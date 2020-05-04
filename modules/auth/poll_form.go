package auth // not sure why this is in package auth?

import (
	"gitea.com/macaron/binding"
	"gitea.com/macaron/macaron"
)

// Form for creating a poll
type CreatePollForm struct {
	Subject     string `binding:"Required;MaxSize(128)"` // 128 is duplicated in the template
	Description string
}

// Validate validates the form fields
func (f *CreatePollForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

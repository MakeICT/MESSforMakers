package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRegEx is a convenience provided for validating email addresses. Recommended by the W3C.
var EmailRegEx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type formErrors map[string][]string

func (e formErrors) Add(f, m string) {
	e[f] = append(e[f], m)
}

func (e formErrors) Get(f string) string {
	es := e[f]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

// Form will hold the form values to be validated and any Errors generated.
// All validation calls are made on the initialized form.
type Form struct {
	url.Values
	Errors formErrors
}

// NewForm adds the form values to the form and initializes the error map. ParseForm should be called before passing http.Request.Form or PostForm to New()
func NewForm(data url.Values) *Form {
	return &Form{
		data,
		formErrors(map[string][]string{}),
	}
}

// Required takes a list of fields that must not be blank in a valid form.
// An error will be added for each field that is blank
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MinLength takes a field name and an int, will add an error to the form if the value in the field is less than the minimum length
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short, must have at least %d characters", d))
	}
}

// MaxLength takes a field name and an int, will add an error to the form if the value in the field is more than the maximum length
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (Max %d characters)", d))
	}
}

// PermittedValues takes the name of a form field and a list of possible options. Will add an error if the value is not in the list of options.
// Useful for making sure radio buttons are not tampered with.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// MatchPattern takes the name of a form field and a regular expression. Adds an error if the field value does not match the expression.
// Useful for validating email addresses, phone numbers, etc.
func (f *Form) MatchPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// MatchField takes the name of two fields. If the fields do not match an error is added.
// Useful for field-confirm patterns
func (f *Form) MatchField(field, match string) {
	value1 := f.Get(field)
	value2 := f.Get(match)
	if value1 == "" || value2 == "" {
		return
	}
	if value1 != value2 {
		f.Errors.Add(field, "Values must match")
	}
}

// Valid should be called only after calling all the other form validators.
// Will return true if there are no Errors on the form.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

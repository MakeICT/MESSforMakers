package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// EmailRegEx is a convenience provided for validating email addresses. Using regex recommended by the W3C.
var EmailRegEx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// PhoneRegEx is a convenience provided for validating phone numbers
var PhoneRegEx = regexp.MustCompile(`^\(?([0-9]{3})\)?[-. ]?([0-9]{3})[-. ]?([0-9]{4})$`)

type formErrors map[string][]string

func (e formErrors) Add(f, m string) {
	e[f] = append(e[f], m)
}

func (e formErrors) Get(f string) string {
	es := e[f]
	if len(es) == 0 {
		return ""
	}
	return strings.Join(es, "; ")
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

// RequiredIf takes a field and a flag. if the flag is true, then the field is required and an error will be added if it's blank.
// Useful for fields that only require information if a specific option is checked.
func (f *Form) RequiredIf(field string, apply bool) {
	if apply {
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
		f.Errors.Add(match, "Values must match")
	}
}

// Valid should be called only after calling all the other form validators.
// Will return true if there are no Errors on the form.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//IntOK returns an int and an OK flag if the string can be converted to an int between max and min inclusive
func IntOK(val string, min, max int) (int, bool) {
	n, err := strconv.Atoi(val)
	if err != nil || n < min || n > max {
		return 0, false
	}
	return n, true
}

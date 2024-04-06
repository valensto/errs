// Package errs provides a simple error handling mechanism based with slug
// tailored for web applications. It defines custom error types
// for various HTTP status codes and offers a utility function
// to convert these errors into their corresponding HTTP status codes.
// This approach simplifies error handling across the HTTP handlers,
// ensuring consistent responses and facilitating easier debugging.
package errs

import (
	"errors"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

// Params represents a map of string key-value pairs. It is used to associate
// additional information with errors, such as field names and validation messages
// in the context of web requests. Params provides methods to manipulate these
// key-value pairs, making it easier to add, check, or retrieve error-related information.
type Params map[string]string

// NewParams creates and returns a new instance of Params. It is a convenience
// function to initialize Params without needing to manually specify its type.
func NewParams() Params {
	return make(Params)
}

// Add inserts a key-value pair into the Params map. If the Params map has not
// been initialized (nil), Add will initialize it before adding the key-value pair.
// This method ensures that key-value pairs can be safely added to Params even
// when it is in its zero value state.
func (p *Params) Add(key, value string) {
	if *p == nil {
		*p = make(Params)
	}
	(*p)[key] = value
}

// IsNil checks if the Params map is empty. It returns true if the map contains
// no key-value pairs, indicating that no additional information has been associated
// with an error.
func (p *Params) IsNil() bool {
	return len(*p) == 0
}

// Err represents a detailed error type in the errs package. It contains
// a Slug indicating the type of error, an underlying error, optional detailed
// messages, and a Params map to hold additional error information. This structure
// allows for rich error descriptions and easy translation or formatting for end users.
type Err struct {
	slug    Slug
	error   error
	details string
	params  Params
}

// Error returns a string representation of the Err. It combines the underlying
// error message with any additional details provided. This method satisfies
// the error interface, allowing Err to be used like any other error object.
func (e Err) Error() string {
	if e.details == "" {
		return e.error.Error()
	}
	return e.error.Error() + ": " + e.details
}

// New creates a new Err instance based on the given Slug and optional details.
// The Slug represents the specific type of error, while the details provide
// further context or information about the error condition. This function is
// a primary entry point for creating typed errors in applications.
func New(slug Slug, details ...string) Err {
	return Err{
		slug:    slug,
		error:   slug,
		details: strings.Join(details, ": "),
	}
}

// NewFromError creates a new Err instance from a generic error object. It attempts
// to preserve the error type if it is already a typed Err; otherwise, it wraps
// the error in a new Err with the SlugUnknown type. This function facilitates
// error conversion and propagation in applications.
func NewFromError(err error) Err {
	return getErr(err)
}

// NewFromValidator creates a new Err instance from a validation error produced
// by the validator package. It translates validation error messages using the
// provided translator and associates them with their corresponding field names
// in the Params map. This function is particularly useful for handling validation
// errors in a structured and user-friendly manner.
func NewFromValidator(err error, translator ut.Translator) Err {
	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		return Err{
			slug:  SlugInvalid,
			error: err,
		}
	}

	e := Err{
		slug:  SlugInvalid,
		error: err,
	}
	for _, err := range err.(validator.ValidationErrors) {
		msg := err.Error()
		if translator != nil {
			msg = err.Translate(translator)
		}

		field := strings.ToLower(err.Field())
		e.params.Add(field, msg)
	}

	return e
}

// WithError enriches the Err instance with an additional underlying error, allowing
// for error wrapping and chaining. This method is useful for building a detailed
// error trace.
func (e Err) WithError(err error) Err {
	e.error = fmt.Errorf("%w: %w", e.error, err)
	return e
}

// WithDetails appends additional details to the Err, providing more context
// about the error. This method helps in conveying precise error information to the end user.
func (e Err) WithDetails(details string) Err {
	e.details = fmt.Sprintf("%s: %s", e.details, details)
	return e
}

// WithParams merges a map of string key-value pairs into the Err's Params, allowing
// for the association of arbitrary data with the error. This method is useful for
// attaching detailed error metadata, such as field-specific error messages.
func (e Err) WithParams(params map[string]string) Err {
	if e.params == nil {
		e.params = make(Params)
	}

	for k, v := range params {
		e.params[k] = v
	}
	return e
}

// DetailFromErr extracts the detailed error message from a given error, if the error is of
// type Err and contains detailed information. This allows for the retrieval of additional
// error context useful for logging or displaying to an end user. If the error does not contain
// detailed information, a default "unknown error" message is returned.
func DetailFromErr(err error) string {
	var e Err
	if errors.As(err, &e) {
		return e.details
	}

	return "unknown error"
}

// ParamsFromErr extracts the Params map from a given error, if the error is of type Err and
// contains a Params map. This is particularly useful for errors that include field-specific
// validation messages or other metadata that might influence how an error is handled or displayed.
// If the error does not contain a Params map, nil is returned.
func ParamsFromErr(err error) map[string]string {
	var e Err
	if errors.As(err, &e) {
		return e.params
	}

	return nil
}

func getErr(err error) Err {
	var e Err
	if !errors.As(err, &e) {
		return Err{
			slug:  SlugUnknown,
			error: err,
		}
	}

	return e
}

package errs

import "errors"

// Slug is an interface that all error slugs implement, providing a basic Error method.
// It allows for a unified way of handling different types of errors through their slug representations.
type Slug interface {
	Error() string
}

type (
	// slug is the base type for all specific slugs, implementing the Slug interface.
	slug string

	// NotFound represents a 404 Not Found error.
	NotFound slug

	// Duplicate represents a 409 Conflict error, indicating a duplicate resource or action.
	Duplicate slug

	// Invalid represents a 400 Bad Request error, typically used for validation failures.
	Invalid slug

	// Forbidden represents a 403 Forbidden error, indicating lack of permission.
	Forbidden slug

	// Unauthorized represents a 401 Unauthorized error, indicating missing or invalid authentication.
	Unauthorized slug

	// NotImplemented represents a 501 Not Implemented error, for unimplemented functionality.
	NotImplemented slug

	// Internal represents a 500 Internal Server Error, indicating a server-side error.
	Internal slug

	// Unknown represents an unspecified error, used as a fallback.
	Unknown slug
)

// Constants for each slug type, providing clear and concise identifiers for common error conditions.
// You can add more slugs as needed to cover additional error cases in your application based on slug types.
const (
	SlugUnknown        Unknown        = "unknown"
	SlugNotFound       NotFound       = "not-found"
	SlugInvalid        Invalid        = "request-invalid"
	SlugUnauthorized   Unauthorized   = "unauthorized"
	SlugForbidden      Forbidden      = "forbidden"
	SlugDuplicate      Duplicate      = "already-exists"
	SlugNotImplemented NotImplemented = "not-implemented"
	SlugInternal       Internal       = "internal-error"
)

// Error returns the string representation of the NotFound error slug,
// facilitating its use as an error object.
func (s Unknown) Error() string {
	return string(s)
}

func (s NotFound) Error() string {
	return string(s)
}

func (s Invalid) Error() string {
	return string(s)
}

func (s Unauthorized) Error() string {
	return string(s)
}

func (s Forbidden) Error() string {
	return string(s)
}

func (s Duplicate) Error() string {
	return string(s)
}

func (s NotImplemented) Error() string {
	return string(s)
}

func (s Internal) Error() string {
	return string(s)
}

// SlugFromErr extracts the Slug from a given error, if the error is of type Err and contains
// a slug. This function is useful for determining the type of an error when handling it,
// especially in situations where the specific error type influences the application's response.
// If the error does not contain a slug, SlugUnknown is returned as a fallback.
func SlugFromErr(err error) (s Slug) {
	var e Err
	if errors.As(err, &e) {
		return e.slug
	}

	return SlugUnknown
}

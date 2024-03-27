package errs

import (
	"errors"
	"net/http"
)

// ProblemJSONType defines a type for representing the "type" field in a Problem Details object.
// According to RFC 9457, the "type" field is a URI reference that identifies the problem type.
// It provides a means to give clients more information about the error in a machine-readable format.
type ProblemJSONType string

const (
	// BlankType is a predefined ProblemJSONType representing a general problem type as specified by
	// RFC 9457. It indicates that the problem does not match any other predefined problem types.
	BlankType ProblemJSONType = "about:blank"
)

// HTTPStatus takes an error object as input and returns the corresponding
// HTTP status code. It works by unwrapping the provided error to its base type
// and then matching it against a set of predefined errors. Each predefined error
// is associated with a specific HTTP status code, allowing for clear and concise
// error handling in web applications. If the error does not match any of the predefined
// types, HTTPStatus defaults to returning http.StatusInternalServerError, indicating
// an unexpected condition.
//
// Usage:
//
//	err := someOperation()
//	if err != nil {
//	    statusCode := errs.HTTPStatus(err)
//	    // Use statusCode for setting HTTP response status
//	}
//
// This function is essential for converting internal error types into
// appropriate HTTP responses, thereby encapsulating the error handling logic
// and promoting a cleaner and more maintainable codebase.
func HTTPStatus(err error) int {
	e := getErr(err)

	errorMap := map[error]int{
		(*NotFound)(nil):       http.StatusNotFound,
		(*BadRequest)(nil):     http.StatusBadRequest,
		(*Unauthorized)(nil):   http.StatusUnauthorized,
		(*Forbidden)(nil):      http.StatusForbidden,
		(*Duplicate)(nil):      http.StatusConflict,
		(*NotImplemented)(nil): http.StatusNotImplemented,
	}

	for errType, statusCode := range errorMap {
		if errors.As(e.slug, &errType) {
			return statusCode
		}
	}

	return http.StatusInternalServerError
}

// ProblemJSON constructs a map representing a Problem Details object as specified by RFC 9457.
// This standard provides a way to carry machine-readable details of errors in a HTTP response to
// avoid the need to define new error response formats for HTTP APIs. The function takes an error,
// an instance URI that identifies the specific occurrence of the problem, and a ProblemJSONType
// indicating the type of problem.
//
// The resulting map includes the following fields:
// - type: A URI reference (ProblemJSONType) that identifies the problem type.
// - title: A short, human-readable summary of the problem type represented by the slug extracted from the error.
// - status: The HTTP status code generated from the error.
// - detail: A human-readable explanation specific to this occurrence of the problem.
// - instance: A URI reference that identifies the specific occurrence of the problem.
// Additional fields, like 'params', provide further details about the problem when available.
//
// Usage:
//
//	err := someOperation()
//	problemDetails := errs.ProblemJSON(err, "http://example.com/err/1234", errs.BlankType)
//	// Use problemDetails to construct the HTTP response body.
func ProblemJSON(err error, instance string, pbType ...ProblemJSONType) map[string]any {
	code := HTTPStatus(err)
	slug := SlugFromErr(err)
	detail := DetailFromErr(err)
	params := ParamsFromErr(err)

	t := BlankType
	if len(pbType) > 0 {
		t = pbType[0]
	}

	errMap := map[string]any{
		"type":     t,
		"title":    slug,
		"status":   code,
		"instance": instance,
	}
	if detail != "" {
		errMap["detail"] = detail
	}

	if len(params) > 0 {
		errMap["params"] = params
	}

	return errMap
}

# errs: Enhanced Error Handling for Go
errs is a Go package that provides enhanced error handling capabilities for Go applications, especially suited for web services and APIs. It introduces structured errors with customizable slugs, HTTP status code mapping, detailed error messages, and support for Problem Details for HTTP APIs (RFC 7807).

## Features
Custom Error Types: Define errors with unique slugs for easy identification.
HTTP Status Code Mapping: Automatically convert errors into appropriate HTTP status codes.
Detailed Error Information: Attach detailed error messages and custom parameters to errors.
Problem Details Support: Generate error responses conforming to RFC 9457 for RESTful APIs.

## Installation
To install errs, use the go get command:

```shell
go get github.com/valensto/errs
```

Ensure you have Go installed and your workspace is properly set up.

Usage
Here's a simple example to get started with errs:

```go
package main

import (
    "fmt"
    "github.com/valensto/errs"
    "net/http"
)

func main() {
    // Create a new error
    err := errs.New(errs.SlugNotFound, "User not found")

    // Convert the error into an HTTP status code
    statusCode := errs.HTTPStatus(err)
    fmt.Printf("HTTP Status Code: %d\n", statusCode)
    
    // Create a Problem Details response
    problemDetails := errs.ProblemJSON(err, "/path/where/error/occurred", errs.BlankType)
    fmt.Printf("Problem Details: %+v\n", problemDetails)
}
```

## Good practices
Every package that can return errors should create a errors.go file looking like this :

```go
package xxx

import "github.com/valensto/errs"

var (
	SlugPhoneInvalid errs.BadRequest = "phone_invalid"
	SlugEmailInvalid errs.BadRequest = "email_invalid"
)
```

In the package, you should always return one of these errors.

```go
// DO THIS
if err != nil {
	return fmt.Errorf("add more details to the trace: %w", err)
}
if !ok {
	return errs.New(SlugPhoneInvalid, "phone is invalid")
}

// DON'T DO THIS
if err != nil {
	return err
}
if !ok {
	return fmt.Errorf("undefined error")
}
if !ok {
	return errors.New("undefined error")
}
```

This way, the user of your package know what to expect regarding errors.

Note: you can always add details with fmt.Errorf("%w: details", ErrFirst)

## Benefits
In your package tests, you can catch every case using errors.Is(err, ErrFirst) or require.ErrorIs(t, err, ErrFirst). You don't have (and shouldn't have) to rely on errors of another package in your tests.

The caller of your package can define different behavior depending on the errors returned.

The root error allows you to trace the execution that lead to the error e.g. invalid request: invalid: wrong format: hex contains invalid characters

## Rule of thumb
- never use errors.New()
- never directly return an error that is not exposed by your package
- wrap errors using fmt.Errorf() with the %w syntax
- use the syntax err: subErr for consistency in logs
- use this package errors at the transport layer ; to find the HTTP/GRPC status and return a known slug to the front-end

Documentation
For more detailed documentation, refer to the code comments within the package.

Contributing
Contributions to errs are welcome!

Support
If you have any questions or issues, please open an issue on the GitHub repository, or contact me directly at hi@valensto.com
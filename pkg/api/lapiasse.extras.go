package api

import (
	"encoding/json/v2"
	"fmt"
)

const unexpectedResponseStructNil = "status code is %d but parsed body %q is nil"

// Make JsonError a proper error type

func (e Error) Error() string {
	if e.Detail == nil {
		return e.Title
	}

	return fmt.Sprintf("%s: %s", e.Title, *e.Detail)
}

func unmarshalJsonErrorOrNil(data []byte) *Error {
	var jsonErr Error
	if err := json.Unmarshal(data, &jsonErr); err != nil {
		return nil
	}

	return &jsonErr
}

// ErrUnexpectedResponseStatusCode is returned when the response has a status
// code that was not expected (not documented).
type ErrUnexpectedResponseStatusCode struct {
	StatusCode int
	JsonError  *Error
}

func (e ErrUnexpectedResponseStatusCode) Error() string {
	return fmt.Sprintf("unexpected response status code: %d", e.StatusCode)
}

// Switch calls the function corresponding to the response's status code.
// It offers an exaustive handling of all documented response status codes.
func (r CreateBookResponse) Switch(
	case201 func(*BookShow) error,
	case400 func(*BadRequestError) error,
	case422 func(*ValidationError) error,
) error {
	switch r.StatusCode() {
	case 201:
		if r.JSON201 == nil {
			return fmt.Errorf(unexpectedResponseStructNil, r.StatusCode(), "JSON201")
		}
		return case201(r.JSON201)
	case 400:
		if r.JSON400 == nil {
			return fmt.Errorf(unexpectedResponseStructNil, r.StatusCode(), "JSON400")
		}
		return case400(r.JSON400)
	case 422:
		if r.JSON422 == nil {
			return fmt.Errorf(unexpectedResponseStructNil, r.StatusCode(), "JSON422")
		}
		return case422(r.JSON422)
	default:
		return ErrUnexpectedResponseStatusCode{
			StatusCode: r.StatusCode(),
			JsonError:  unmarshalJsonErrorOrNil(r.Body),
		}
	}
}

// Switch calls the function corresponding to the response's status code.
// It offers an exaustive handling of all documented response status codes.
func (r UpdateBookResponse) Switch(
	case200 func(*BookShow) error,
	case400 func(*BadRequestError) error,
	case404 func(*NotFoundError) error,
	case422 func(*ValidationError) error,
) error {
	switch r.StatusCode() {
	case 200:
		if r.JSON200 == nil {
			return fmt.Errorf(unexpectedResponseStructNil, r.StatusCode(), "JSON200")
		}
		return case200(r.JSON200)
	case 400:
		if r.JSON400 == nil {
			return fmt.Errorf(unexpectedResponseStructNil, r.StatusCode(), "JSON400")
		}
		return case400(r.JSON400)
	case 404:
		if r.JSON404 == nil {
			return fmt.Errorf(unexpectedResponseStructNil, r.StatusCode(), "JSON404")
		}
		return case404(r.JSON404)
	case 422:
		if r.JSON422 == nil {
			return fmt.Errorf(unexpectedResponseStructNil, r.StatusCode(), "JSON422")
		}
		return case422(r.JSON422)
	default:
		return ErrUnexpectedResponseStatusCode{
			StatusCode: r.StatusCode(),
			JsonError:  unmarshalJsonErrorOrNil(r.Body),
		}
	}
}

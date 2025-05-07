package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ellypaws/inkbunny/types"
)

// ParseResponse parses the HTTP response and returns the decoded value of type T.
// It checks [http.Response.StatusCode], decodes and checks if the [http.Response.Body]
// decodes into types.ErrorResponse, and finally decodes into T if no errors are returned.
// ParseResponse also calls [io.Closer.Close] on the Body.
func ParseResponse[T any](response *http.Response) (T, error) {
	var t T
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return t, fmt.Errorf("unexpected status code %s (%d)", response.Status, response.StatusCode)
	}

	bin, err := io.ReadAll(response.Body)
	if err != nil {
		return t, err
	}

	errResponse, err := DecodeBytes[types.ErrorResponse](bin)
	if err != nil {
		return t, err
	}

	if errResponse.Code != nil {
		return t, fmt.Errorf("[%d]: %s", *errResponse.Code, errResponse.Message)
	}

	return DecodeBytes[T](bin)
}

func Decode[T any](body io.Reader) (T, error) {
	d := json.NewDecoder(body)
	var v T
	if err := d.Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}

func DecodeBytes[T any](body []byte) (T, error) {
	return Decode[T](bytes.NewReader(body))
}

func Encode[T any](v T, w io.Writer) error {
	e := json.NewEncoder(w)
	if err := e.Encode(v); err != nil {
		return err
	}
	return nil
}

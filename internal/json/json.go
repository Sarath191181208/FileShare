package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)


func ReadJSON(data interface{}, w http.ResponseWriter, r *http.Request) error {
	// Limiting JSON body to be 1MB
	maxBytes := 1 * 1024 * 1024 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Creating a decoder
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decoding the json data
	err := decoder.Decode(&data)

	// if no error return
	if err == nil {
		// decode an empty struct for double json data ex:
		// '{"k": "v"}  {"k":"v"}'
		err = decoder.Decode(&struct{}{})
		if err != io.EOF {
			return errors.New("body must only contain a single JSON value")
		}
		return nil
	}

	var syntaxError *json.SyntaxError
	var unmarshallTypeError *json.UnmarshalTypeError
	unknownFiledErrorPrefix := "json: unknown field "

	switch {
	// catching JSON with invalid syntax
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly formatted JSON (at character %d)", syntaxError.Offset)

	// catching Non proper JSON
	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("body contains badly formed JSON")

	// catching invalid json type for a files
	// ex: "30" insted of 30
	case errors.As(err, &unmarshallTypeError):
		invalidJsonType := unmarshallTypeError.Field != ""
		if invalidJsonType {
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshallTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshallTypeError.Offset)

	// catching Unkown Filed error
	case strings.HasPrefix(err.Error(), unknownFiledErrorPrefix):
		filedName := strings.TrimPrefix(err.Error(), unknownFiledErrorPrefix)
		return fmt.Errorf("body contains unkonwn key %s", filedName)

	// catching max bytes limit
	case errors.Is(err, &http.MaxBytesError{}):
		return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// catching end of file i.e empty body
	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")

	default:
		return err
	}
}

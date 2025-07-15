package server

import (
	"encoding/json"
	"io"
)

func marshal[T any](w io.Writer, v T) error {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return err
	}

	return nil
}

func unmarshal[T any](r io.Reader) (T, error) {
	var v T

	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		return v, err
	}

	return v, nil
}

package jsons

import (
	"encoding/json"
	"io"
)

type GenericFormatter struct{}

func NewGenericFormatter() *GenericFormatter { return &GenericFormatter{} }

func (x GenericFormatter) Format(w io.Writer, data any) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

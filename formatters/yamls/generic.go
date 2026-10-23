package yamls

import (
	"io"

	"gopkg.in/yaml.v3"
)

type GenericFormatter struct{}

func NewGenericFormatter() *GenericFormatter { return &GenericFormatter{} }

func (x GenericFormatter) Format(w io.Writer, data any) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

package ports

import "io"

type Formatter interface {
	Format(w io.Writer, data any) error
}

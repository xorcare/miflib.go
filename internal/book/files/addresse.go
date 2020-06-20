package files

import (
	"github.com/xorcare/miflib.go/internal/jstring"
)

// Address contains the parameters of the file by which you can download it.
type Address struct {
	URL      string         `json:"url,omitempty"`
	Size     uint           `json:"size,omitempty"`
	Duration string         `json:"duration,omitempty"`
	Title    jstring.String `json:"title,omitempty"`
}

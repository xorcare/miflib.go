package book

import (
	"encoding/json"
)

// Formats it's relation file formats to their addresses.
type Formats map[string]Addresses

func (m Formats) String() string {
	text, err := json.Marshal(m)
	if err != nil {
		panic("files: Formats.String(): " + err.Error())
	}
	return string(text)
}

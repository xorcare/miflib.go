package files

import (
	"encoding/json"
)

// Map it's relation file type to their addresses.
type Map map[string]Addresses

func (m Map) String() string {
	text, err := json.Marshal(m)
	if err != nil {
		panic("files: Map.String(): " + err.Error())
	}
	return string(text)
}

// Files information about all available files for download.
type Files struct {
	Books      Map `json:"ebook"`
	AudioBooks Map `json:"audiobook"`
	Demo       Map `json:"demo"`
}

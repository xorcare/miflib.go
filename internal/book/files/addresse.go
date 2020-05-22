package files

// Address contains the parameters of the file by which you can download it.
type Address struct {
	URL      string `json:"url,omitempty"`
	Size     uint   `json:"size,omitempty"`
	Duration string `json:"duration,omitempty"`
	Title    string `json:"title,omitempty"`
}

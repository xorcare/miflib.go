package files

// Files information about all available files for download.
type Files struct {
	Books      map[string]Addresses `json:"ebook"`
	AudioBooks map[string]Addresses `json:"audiobook"`
	Demo       map[string]Addresses `json:"demo"`
}

package downloader

import (
	"path/filepath"
	"strings"

	"github.com/xorcare/miflib.go/internal/norm"
)

// https://en.wikipedia.org/wiki/Filename
// https://support.microsoft.com/en-us/office/invalid-file-names-and-file-types-in-onedrive-and-sharepoint-64883a5d-228e-48f5-b3d2-eb39e07630fa
// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file

var forbiddenChars = unique([]string{
	// FAT12, FAT16, FAT32
	`!`, // exclamation mark
	`"`, // double quote
	`*`, // asterisk
	`/`, // forward slash
	`:`, // colon
	`<`, // less than
	`>`, // greater than
	`?`, // question mark
	`@`, // at
	`\`, // backslash
	`|`, // vertical bar or pipe

	// exFAT, NTFS, VFAT
	`"`, // double quote
	`*`, // asterisk
	`/`, // forward slash
	`:`, // colon
	`<`, // less than
	`>`, // greater than
	`?`, // question mark
	`\`, // backslash
	`|`, // vertical bar or pipe

	// Windows
	`"`, // double quote
	`%`, // percent
	`*`, // asterisk
	`/`, // forward slash
	`:`, // colon
	`<`, // less than
	`>`, // greater than
	`?`, // question mark
	`\`, // backslash
	`|`, // vertical bar or pipe

	// Mac OS HFS,HFS+
	`:`, // colon

	// OneDrive
	`"`, // double quote
	`*`, // asterisk
	`/`, // forward slash
	`:`, // colon
	`<`, // less than
	`>`, // greater than
	`?`, // question mark
	`\`, // backslash
	`|`, // vertical bar or pipe
	`\`, // backslash

	// OneDrive for business, SharePoint Server 2013
	`~`, // swung dash or tilde
	`"`, // double quote
	`*`, // asterisk
	`/`, // forward slash
	`:`, // colon
	`<`, // less than
	`>`, // greater than
	`?`, // question mark
	`\`, // backslash
	`{`, // opening braces
	`|`, // vertical bar or pipe
	`}`, // closing curly brackets
})

// replacer it's replaces the specified characters with nothing.
var replacer = strings.NewReplacer(append(strings.Split(strings.Join(forbiddenChars, "  "), " "), "")...)

func clearBaseName(s string) string {
	old := s
	s = norm.String(s)
	s = replacer.Replace(s)
	s = strings.ReplaceAll(s, " .", ".")
	s = strings.TrimSuffix(s, ".")
	s = norm.String(s)

	if old == s {
		return s
	}

	return clearBaseName(s)
}

func clearBase(s string) string {
	return strings.ReplaceAll(s, filepath.Base(s), clearBaseName(filepath.Base(s)))
}

func unique(s []string) []string {
	keys := make(map[string]bool, len(s))
	var list []string
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

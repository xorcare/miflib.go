package downloader

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xorcare/golden"
)

func Test_clearBaseName(t *testing.T) {
	tests := map[string]string{
		// write specific test cases.
		"\n": "",   // line feed
		"\r": "",   // carriage return
		"\t": "",   // horizontal tab
		`!`:  "",   // exclamation mark
		`"`:  "",   // double quote
		`%`:  "",   // percent
		`*`:  "",   // asterisk
		`/`:  "",   // forward slash
		`:`:  "",   // colon
		`<`:  "",   // less than
		`>`:  "",   // greater than
		`?`:  "",   // question mark
		`@`:  "",   // at
		`\`:  "",   // backslash
		`{`:  "",   // opening braces
		`|`:  "",   // vertical bar or pipe
		`}`:  "",   // closing curly brackets
		`~`:  "",   // swung dash or tilde
		`№`:  "No", // number sign

	}

	// load big set of test cases.
	read := golden.Read(t)
	require.NotEmpty(t, read)
	require.NoError(t, json.Unmarshal(read, &tests))
	for arg, want := range tests {
		t.Run(arg, func(t *testing.T) {
			got := clearBaseName(arg)
			require.Equal(t, want, got)
		})
	}
}

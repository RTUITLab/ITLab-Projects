package landingparser

import (
	"bytes"
)

func PrepareLandingToParse(
	data []byte,
) []byte {
	return bytes.Replace(
		bytes.ReplaceAll(
			data,
			[]byte("\r"),
			[]byte(""),
		),
		[]byte("---\n"),
		[]byte("\n---\n\n"),
		-1,
	)
}
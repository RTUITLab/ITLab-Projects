package landingparser

import (
	"bytes"
)

func PrepareLandingToParse(
	data []byte,
) []byte {
	return bytes.Replace(
		data,
		[]byte("---\n"),
		[]byte("\n---\n\n"),
		-1,
	)
}
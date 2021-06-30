package urlvalue

import (
	"net/url"
	"strconv"
	"strings"
)

func ParseStringsFromURL(
	Map		MapValuer,
	Values	url.Values,
	keys	...string,
) {
	for _, k := range keys {
		Map.SetString(
			k,
			ParseMassOfStringsToString(
				Values[k],
			),
		)
	}
}

func ParseNumbersFromURL(
	Map		MapValuer,
	Values	url.Values,
	keys	...string,
) {
	for _, k := range keys {
		Map.SetInt(
			k,
			int(
				ParseStringToInt(
					ParseMassOfStringsToString(
						Values[k],
					),
				),
			),
		)
	}
}

func ParseStringToInt(
	value	string,
) int64 {
	v, err := strconv.ParseInt(
		value,
		10,
		64,
	)
	if err != nil {
		return 0
	}

	return v
}

func ParseMassOfStringsToString(
	mass []string,
) string {
	builder := strings.Builder{}

	for i, s := range mass {
		builder.WriteString(s)
		if !(i == len(mass) - 1) {
			builder.WriteString(" ")
		}
	}

	return builder.String()
}
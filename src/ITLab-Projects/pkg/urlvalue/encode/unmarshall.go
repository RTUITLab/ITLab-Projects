package encode

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/ITLab-Projects/pkg/urlvalue"
	"github.com/pkg/errors"
)

var (
	TypeError 	= errors.New("Type error")
	ValueError	= errors.New("Value is not pointer")
)

type UnmarshallOptions struct {
	Strict bool
}

type UrlQueryDecoder struct {
	Opts	UnmarshallOptions
}

func NewUrlQueryDecode() *UrlQueryDecoder {
	return &UrlQueryDecoder{
		Opts: UnmarshallOptions{
			Strict: false,
		},
	}
}

type UrlQueryUnmarshaler interface {
	UnmarshalUrlQuery(
		values	url.Values,
	) error
}

func (d *UrlQueryDecoder) UrlQueryUnmarshall(
	v		interface{},
	values	url.Values,
) error {
	u, ok := v.(UrlQueryUnmarshaler)
	if ok {
		return u.UnmarshalUrlQuery(
			values,
		)
	}
	value := reflect.ValueOf(v)
	t := value.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		return ValueError
	}
	value = value.Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if fieldIsStruct(field) {
			vField := value.Field(i)
			if vField.Kind() != reflect.Ptr {
				d.UrlQueryUnmarshall(
					vField.Addr().Interface(),
					values,
				)
			} else {
				d.UrlQueryUnmarshall(
					vField.Interface(),
					values,
				)
			}
		}
		opts, ok := field.Tag.Lookup("query")
		if !ok {
			continue
		}
		err := unmarshallFieldFromURL(
			value.Field(i),
			opts,
			values,
		)
		if err := d.handleError(err); err != nil {
			return err
		}
	}
	return nil
}

func (d *UrlQueryDecoder) handleError(
	err error,
) error {
	if err == nil {
		return nil
	} else if d.Opts.Strict {
		return err
	} else {
		switch {
		case 	errors.Is(err, TypeError):
				//   Pass
		default:
			return err
		}
	}
	return nil
}

func UrlQueryUnmarshall(
	v		interface{},
	values	url.Values,
) error {
	return NewUrlQueryDecode().UrlQueryUnmarshall(v, values)
}

type queryValue struct {
	key		string
	_type	queryTypes
}

type queryTypes int

const (
	STRING queryTypes	= iota
	INT
)

func parseOpts(stropts string) *queryValue {
	opts := strings.Split(
		strings.ReplaceAll(
			stropts,
			" ",
			"",
		),
		",",
	)

	if len(opts) == 0 {
		return nil
	}

	if len(opts) == 1 {
		return &queryValue{
			key: opts[0],
			_type: STRING,
		}
	}

	opt := &queryValue{
		key: opts[0],
	}
	switch opts[1] {
	case "int", "integer":
		opt._type = INT
	case "string":
		opt._type = STRING
	}

	return opt
}

func unmarshallFieldFromURL(
	field 		reflect.Value, 
	optsString 	string,
	values		url.Values,
) error {
	opts := parseOpts(optsString)

	return unmarshallField(
		field,
		opts._type,
		urlvalue.ParseMassOfStringsToString(
			values[opts.key],
		),
	)
}

func unmarshallField(
	field 		reflect.Value, 
	_type 		queryTypes,
	value		string,
) error {
	// TODO make avaliable to scan to pointer
	switch _type {
	case STRING:
		field.SetString(value)
	case INT:
		intValue, err := strconv.ParseInt(
			value,
			10,
			64,
		)
		if err != nil {
			return errors.Wrap(
				TypeError,
				fmt.Sprintf(
					"Error on unmarshall field %s: %v",
					field.Type().Name(), err,
				),
			)
		}

		field.SetInt(intValue)
	}
	return nil
}

func fieldIsStruct(
	field reflect.StructField,
) bool {
	if field.Type.Kind() == reflect.Struct {
		return true
	} else if field.Type.Kind() == reflect.Ptr {
		return field.Type.Elem().Kind() == reflect.Struct
	}

	return false
}
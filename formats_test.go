package validate

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/circl-dev/spec"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

// Validator for string formats
func TestFormatValidator_EdgeCases(t *testing.T) {
	// Apply
	v := formatValidator{
		KnownFormats: strfmt.Default,
	}

	// formatValidator applies to: Items, Parameter,Schema

	p := spec.Parameter{}
	p.Typed(stringType, "email")
	s := spec.Schema{}
	s.Typed(stringType, "uuid")
	i := spec.Items{}
	i.Typed(stringType, "datetime")

	sources := []interface{}{&p, &s, &i}

	for _, source := range sources {
		// Default formats for strings
		assert.True(t, v.Applies(source, reflect.String))
		// Do not apply for number formats
		assert.False(t, v.Applies(source, reflect.Int))
	}

	assert.False(t, v.Applies("A string", reflect.String))
	assert.False(t, v.Applies(nil, reflect.String))
}

func TestStringValidation(t *testing.T) {
	type testParams struct {
		format string
		obj    fmt.Stringer
	}

	testCases := []*testParams{
		&testParams{
			format: "datetime",
			obj:    strfmt.NewDateTime(),
		},
		&testParams{
			format: "uuid",
			obj:    strfmt.UUID("00000000-0000-0000-0000-000000000000"),
		},
		&testParams{
			format: "email",
			obj:    strfmt.Email("name@domain.tld"),
		},
		&testParams{
			format: "bsonobjectid",
			obj:    strfmt.NewObjectId("60a7903427a1e6666d2b998c"),
		},
	}

	for _, v := range testCases {
		err := FormatOf("id", "body", v.format, v.obj.String(), strfmt.Default)
		assert.Nil(t, err)
	}
}

package jsonprovider

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/configs/configschema"
)

func TestMarshalAttribute(t *testing.T) {
	tests := []struct {
		Input          *configschema.Attribute
		SensitivePaths *configschema.SensitivePathElement
		Want           *attribute
	}{
		{
			&configschema.Attribute{Type: cty.String, Optional: true, Computed: true},
			configschema.NewSensitivePathLeaf(false),
			&attribute{
				AttributeType: json.RawMessage(`"string"`),
				Optional:      true,
				Computed:      true,
			},
		},
		{ // collection types look a little odd.
			&configschema.Attribute{Type: cty.Map(cty.String), Optional: true, Computed: true},
			configschema.NewSensitivePathLeaf(false),
			&attribute{
				AttributeType: json.RawMessage(`["map","string"]`),
				Optional:      true,
				Computed:      true,
			},
		},
		{
			&configschema.Attribute{Type: cty.String, Optional: true, Computed: true},
			configschema.NewSensitivePathLeaf(true),
			&attribute{
				AttributeType: json.RawMessage(`"string"`),
				Optional:      true,
				Computed:      true,
				Sensitive:     true,
			},
		},
		{
			&configschema.Attribute{Type: cty.String, Optional: true, Computed: true},
			&configschema.SensitivePathElement{
				NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
					"anySensitiveChild": configschema.NewSensitivePathLeaf(true),
				},
			},
			&attribute{
				AttributeType: json.RawMessage(`"string"`),
				Optional:      true,
				Computed:      true,
				Sensitive:     false,
			},
		},
	}

	for _, test := range tests {
		got := marshalAttribute(test.Input, test.SensitivePaths)
		if !cmp.Equal(got, test.Want) {
			t.Fatalf("wrong result:\n %v\n", cmp.Diff(got, test.Want))
		}
	}
}

package jsonprovider

import (
	"encoding/json"

	"github.com/hashicorp/terraform/configs/configschema"
)

type attribute struct {
	AttributeType json.RawMessage `json:"type,omitempty"`
	Description   string          `json:"description,omitempty"`
	Required      bool            `json:"required,omitempty"`
	Optional      bool            `json:"optional,omitempty"`
	Computed      bool            `json:"computed,omitempty"`
	Sensitive     bool            `json:"sensitive,omitempty"`
}

func marshalAttribute(attr *configschema.Attribute, sensitivePaths *configschema.SensitivePathElement) *attribute {
	// we're not concerned about errors because at this point the schema has
	// already been checked and re-checked.
	attrTy, _ := attr.Type.MarshalJSON()

	return &attribute{
		AttributeType: attrTy,
		Description:   attr.Description,
		Required:      attr.Required,
		Optional:      attr.Optional,
		Computed:      attr.Computed,
		// We are not marshalling any information about further nested sensitive values.
		// Maybe we should deprecate the sensitive field in this attribute struct
		// and marshal the whole sensitivePaths structure.
		Sensitive: sensitivePaths.IsSensitive(),
	}
}

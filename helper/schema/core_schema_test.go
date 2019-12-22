package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/configs/configschema"
)

// add the implicit "id" attribute for test resources
func testResource(block *configschema.Block) *configschema.Block {
	if block.Attributes == nil {
		block.Attributes = make(map[string]*configschema.Attribute)
	}

	if block.BlockTypes == nil {
		block.BlockTypes = make(map[string]*configschema.NestedBlock)
	}

	if block.Attributes["id"] == nil {
		block.Attributes["id"] = &configschema.Attribute{
			Type:     cty.String,
			Optional: true,
			Computed: true,
		}
	}
	return block
}

func TestSchemaMapCoreConfigSchema(t *testing.T) {
	tests := map[string]struct {
		Schema map[string]*Schema
		Want   *configschema.Block
	}{
		//"empty": {
		//	map[string]*Schema{},
		//	testResource(&configschema.Block{}),
		//},
		//"primitives": {
		//	map[string]*Schema{
		//		"int": {
		//			Type:        TypeInt,
		//			Required:    true,
		//			Description: "foo bar baz",
		//		},
		//		"float": {
		//			Type:     TypeFloat,
		//			Optional: true,
		//		},
		//		"bool": {
		//			Type:     TypeBool,
		//			Computed: true,
		//		},
		//		"string": {
		//			Type:     TypeString,
		//			Optional: true,
		//			Computed: true,
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{
		//			"int": {
		//				Type:        cty.Number,
		//				Required:    true,
		//				Description: "foo bar baz",
		//			},
		//			"float": {
		//				Type:     cty.Number,
		//				Optional: true,
		//			},
		//			"bool": {
		//				Type:     cty.Bool,
		//				Computed: true,
		//			},
		//			"string": {
		//				Type:     cty.String,
		//				Optional: true,
		//				Computed: true,
		//			},
		//		},
		//		BlockTypes: map[string]*configschema.NestedBlock{},
		//	}),
		//},
		//"simple collections": {
		//	map[string]*Schema{
		//		"list": {
		//			Type:     TypeList,
		//			Required: true,
		//			Elem: &Schema{
		//				Type: TypeInt,
		//			},
		//		},
		//		"set": {
		//			Type:     TypeSet,
		//			Optional: true,
		//			Elem: &Schema{
		//				Type: TypeString,
		//			},
		//		},
		//		"map": {
		//			Type:     TypeMap,
		//			Optional: true,
		//			Elem: &Schema{
		//				Type: TypeBool,
		//			},
		//		},
		//		"map_default_type": {
		//			Type:     TypeMap,
		//			Optional: true,
		//			// Maps historically don't have elements because we
		//			// assumed they would be strings, so this needs to work
		//			// for pre-existing schemas.
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{
		//			"list": {
		//				Type:     cty.List(cty.Number),
		//				Required: true,
		//			},
		//			"set": {
		//				Type:     cty.Set(cty.String),
		//				Optional: true,
		//			},
		//			"map": {
		//				Type:     cty.Map(cty.Bool),
		//				Optional: true,
		//			},
		//			"map_default_type": {
		//				Type:     cty.Map(cty.String),
		//				Optional: true,
		//			},
		//		},
		//		BlockTypes: map[string]*configschema.NestedBlock{},
		//	}),
		//},
		//"incorrectly-specified collections": {
		//	// Historically we tolerated setting a type directly as the Elem
		//	// attribute, rather than a Schema object. This is common enough
		//	// in existing provider code that we must support it as an alias
		//	// for a schema object with the given type.
		//	map[string]*Schema{
		//		"list": {
		//			Type:     TypeList,
		//			Required: true,
		//			Elem:     TypeInt,
		//		},
		//		"set": {
		//			Type:     TypeSet,
		//			Optional: true,
		//			Elem:     TypeString,
		//		},
		//		"map": {
		//			Type:     TypeMap,
		//			Optional: true,
		//			Elem:     TypeBool,
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{
		//			"list": {
		//				Type:     cty.List(cty.Number),
		//				Required: true,
		//			},
		//			"set": {
		//				Type:     cty.Set(cty.String),
		//				Optional: true,
		//			},
		//			"map": {
		//				Type:     cty.Map(cty.Bool),
		//				Optional: true,
		//			},
		//		},
		//		BlockTypes: map[string]*configschema.NestedBlock{},
		//	}),
		//},
		//"sub-resource collections": {
		//	map[string]*Schema{
		//		"list": {
		//			Type:     TypeList,
		//			Required: true,
		//			Elem: &Resource{
		//				Schema: map[string]*Schema{},
		//			},
		//			MinItems: 1,
		//			MaxItems: 2,
		//		},
		//		"set": {
		//			Type:     TypeSet,
		//			Required: true,
		//			Elem: &Resource{
		//				Schema: map[string]*Schema{},
		//			},
		//		},
		//		"map": {
		//			Type:     TypeMap,
		//			Optional: true,
		//			Elem: &Resource{
		//				Schema: map[string]*Schema{},
		//			},
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{
		//			// This one becomes a string attribute because helper/schema
		//			// doesn't actually support maps of resource. The given
		//			// "Elem" is just ignored entirely here, which is important
		//			// because that is also true of the helper/schema logic and
		//			// existing providers rely on this being ignored for
		//			// correct operation.
		//			"map": {
		//				Type:     cty.Map(cty.String),
		//				Optional: true,
		//			},
		//		},
		//		BlockTypes: map[string]*configschema.NestedBlock{
		//			"list": {
		//				Nesting:  configschema.NestingList,
		//				Block:    configschema.Block{},
		//				MinItems: 1,
		//				MaxItems: 2,
		//			},
		//			"set": {
		//				Nesting:  configschema.NestingSet,
		//				Block:    configschema.Block{},
		//				MinItems: 1, // because schema is Required
		//			},
		//		},
		//	}),
		//},
		//"sub-resource collections minitems+optional": {
		//	// This particular case is an odd one where the provider gives
		//	// conflicting information about whether a sub-resource is required,
		//	// by marking it as optional but also requiring one item.
		//	// Historically the optional-ness "won" here, and so we must
		//	// honor that for compatibility with providers that relied on this
		//	// undocumented interaction.
		//	map[string]*Schema{
		//		"list": {
		//			Type:     TypeList,
		//			Optional: true,
		//			Elem: &Resource{
		//				Schema: map[string]*Schema{},
		//			},
		//			MinItems: 1,
		//			MaxItems: 1,
		//		},
		//		"set": {
		//			Type:     TypeSet,
		//			Optional: true,
		//			Elem: &Resource{
		//				Schema: map[string]*Schema{},
		//			},
		//			MinItems: 1,
		//			MaxItems: 1,
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{},
		//		BlockTypes: map[string]*configschema.NestedBlock{
		//			"list": {
		//				Nesting:  configschema.NestingList,
		//				Block:    configschema.Block{},
		//				MinItems: 0,
		//				MaxItems: 1,
		//			},
		//			"set": {
		//				Nesting:  configschema.NestingSet,
		//				Block:    configschema.Block{},
		//				MinItems: 0,
		//				MaxItems: 1,
		//			},
		//		},
		//	}),
		//},
		//"sub-resource collections minitems+computed": {
		//	map[string]*Schema{
		//		"list": {
		//			Type:     TypeList,
		//			Computed: true,
		//			Elem: &Resource{
		//				Schema: map[string]*Schema{},
		//			},
		//			MinItems: 1,
		//			MaxItems: 1,
		//		},
		//		"set": {
		//			Type:     TypeSet,
		//			Computed: true,
		//			Elem: &Resource{
		//				Schema: map[string]*Schema{},
		//			},
		//			MinItems: 1,
		//			MaxItems: 1,
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{
		//			"list": {
		//				Type:     cty.List(cty.EmptyObject),
		//				Computed: true,
		//			},
		//			"set": {
		//				Type:     cty.Set(cty.EmptyObject),
		//				Computed: true,
		//			},
		//		},
		//	}),
		//},
		//"nested attributes and blocks": {
		//	map[string]*Schema{
		//		"foo": {
		//			Type:     TypeList,
		//			Required: true,
		//			Elem: &Resource{
		//				Schema: map[string]*Schema{
		//					"bar": {
		//						Type:     TypeList,
		//						Required: true,
		//						Elem: &Schema{
		//							Type: TypeList,
		//							Elem: &Schema{
		//								Type: TypeString,
		//							},
		//						},
		//					},
		//					"baz": {
		//						Type:     TypeSet,
		//						Optional: true,
		//						Elem: &Resource{
		//							Schema: map[string]*Schema{},
		//						},
		//					},
		//				},
		//			},
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{},
		//		BlockTypes: map[string]*configschema.NestedBlock{
		//			"foo": &configschema.NestedBlock{
		//				Nesting: configschema.NestingList,
		//				Block: configschema.Block{
		//					Attributes: map[string]*configschema.Attribute{
		//						"bar": {
		//							Type:     cty.List(cty.List(cty.String)),
		//							Required: true,
		//						},
		//					},
		//					BlockTypes: map[string]*configschema.NestedBlock{
		//						"baz": {
		//							Nesting: configschema.NestingSet,
		//							Block:   configschema.Block{},
		//						},
		//					},
		//				},
		//				MinItems: 1, // because schema is Required
		//			},
		//		},
		//	}),
		//},
		"sensitive attributes": {
			map[string]*Schema{
				"string": {
					Type:      TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"bool": {
					Type:      TypeBool,
					Optional:  true,
					Sensitive: true,
				},
				"int": {
					Type:      TypeInt,
					Optional:  true,
					Sensitive: true,
				},
				"float": {
					Type:      TypeFloat,
					Optional:  true,
					Sensitive: true,
				},
				"listWithSchemaElement": {
					Type:      TypeList,
					Optional:  true,
					Sensitive: true,
					Elem: &Schema{
						Type:     TypeInt,
						Optional: true,
					},
				},
				"listWithSensitiveSchemaElement": {
					Type:       TypeList,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Elem: &Schema{
						Type:      TypeInt,
						Optional:  true,
						Sensitive: true,
					},
				},
				"listWithResourceElement": {
					Type:       TypeList,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"myint": {
								Type:     TypeInt,
								Optional: true,
							},
						},
					},
				},
				"listWithSensitiveResourceElement": {
					Type:       TypeList,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"myint": {
								Type:      TypeInt,
								Optional:  true,
								Sensitive: true,
							},
						},
					},
				},
				"listWithValueTypeElement": {
					Type:       TypeList,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
					Elem:       TypeInt,
				},
				"listWithoutElementDefinition": {
					Type:       TypeList,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
				},
				"setWithSchemaElement": {
					Type:      TypeSet,
					Optional:  true,
					Sensitive: true,
					Elem: &Schema{
						Type:     TypeInt,
						Optional: true,
					},
				},
				"setWithSensitiveSchemaElement": {
					Type:       TypeSet,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Elem: &Schema{
						Type:      TypeInt,
						Optional:  true,
						Sensitive: true,
					},
				},
				"setWithResourceElement": {
					Type:       TypeSet,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"myint": {
								Type:     TypeInt,
								Optional: true,
							},
						},
					},
				},
				"setWithSensitiveResourceElement": {
					Type:       TypeSet,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"myint": {
								Type:      TypeInt,
								Optional:  true,
								Sensitive: true,
							},
						},
					},
				},
				"setWithValueTypeElement": {
					Type:       TypeSet,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
					Elem:       TypeInt,
				},
				"setWithoutElementDefinition": {
					Type:       TypeSet,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
				},
				"mapWithSchemaElement": {
					Type:      TypeMap,
					Optional:  true,
					Sensitive: true,
					Elem: &Schema{
						Type:     TypeInt,
						Optional: true,
					},
				},
				"mapWithSensitiveSchemaElement": {
					Type:       TypeMap,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Elem: &Schema{
						Type:      TypeInt,
						Optional:  true,
						Sensitive: true,
					},
				},
				"unsupported_mapWithResourceElement": {
					Type:       TypeMap,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
					// This is not supported and therefore converted into a schema of type string
					Elem: &Resource{
						Schema: map[string]*Schema{
							"myint": {
								Type:     TypeInt,
								Optional: true,
							},
						},
					},
				},
				"unsupported_mapWithSensitiveResourceElement": {
					Type:       TypeMap,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					// This is not supported and therefore converted into a non sensitive Schema of type string
					Elem: &Resource{
						Schema: map[string]*Schema{
							"myint": {
								Type:      TypeInt,
								Optional:  true,
								Sensitive: true,
							},
						},
					},
				},
				"mapWithValueTypeElement": {
					Type:       TypeMap,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
					Elem:       TypeInt,
				},
				"mapWithoutElementDefinition": {
					Type:       TypeMap,
					ConfigMode: SchemaConfigModeAttr,
					Optional:   true,
					Sensitive:  true,
				},
			},
			testResource(&configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"string": {
						Type:     cty.String,
						Optional: true,
					},
					"bool": {
						Type:     cty.Bool,
						Optional: true,
					},
					"int": {
						Type:     cty.Number,
						Optional: true,
					},
					"float": {
						Type:     cty.Number,
						Optional: true,
					},
					"listWithSchemaElement": {
						Type:     cty.List(cty.Number),
						Optional: true,
					},
					"listWithSensitiveSchemaElement": {
						Type:     cty.List(cty.Number),
						Optional: true,
					},
					"listWithResourceElement": {
						Type: cty.List(cty.Object(map[string]cty.Type{
							"myint": cty.Number,
						})),
						Optional: true,
					},
					"listWithSensitiveResourceElement": {
						Type: cty.List(cty.Object(map[string]cty.Type{
							"myint": cty.Number,
						})),
						Optional: true,
					},
					"listWithValueTypeElement": {
						Type:     cty.List(cty.Number),
						Optional: true,
					},
					"listWithoutElementDefinition": {
						Type:     cty.List(cty.String),
						Optional: true,
					},
					"setWithSchemaElement": {
						Type:     cty.Set(cty.Number),
						Optional: true,
					},
					"setWithSensitiveSchemaElement": {
						Type:     cty.Set(cty.Number),
						Optional: true,
					},
					"setWithResourceElement": {
						Type: cty.Set(cty.Object(map[string]cty.Type{
							"myint": cty.Number,
						})),
						Optional: true,
					},
					"setWithSensitiveResourceElement": {
						Type: cty.Set(cty.Object(map[string]cty.Type{
							"myint": cty.Number,
						})),
						Optional: true,
					},
					"setWithValueTypeElement": {
						Type:     cty.Set(cty.Number),
						Optional: true,
					},
					"setWithoutElementDefinition": {
						Type:     cty.Set(cty.String),
						Optional: true,
					},
					"mapWithSchemaElement": {
						Type:     cty.Map(cty.Number),
						Optional: true,
					},
					"mapWithSensitiveSchemaElement": {
						Type:     cty.Map(cty.Number),
						Optional: true,
					},
					"unsupported_mapWithResourceElement": {
						Type:     cty.Map(cty.String),
						Optional: true,
					},
					"unsupported_mapWithSensitiveResourceElement": {
						Type:     cty.Map(cty.String),
						Optional: true,
					},
					"mapWithValueTypeElement": {
						Type:     cty.Map(cty.Number),
						Optional: true,
					},
					"mapWithoutElementDefinition": {
						Type:     cty.Map(cty.String),
						Optional: true,
					},
				},
				BlockTypes: map[string]*configschema.NestedBlock{},
				SensitivePaths: &configschema.SensitivePathElement{
					NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
						"string":                {},
						"bool":                  {},
						"int":                   {},
						"float":                 {},
						"listWithSchemaElement": {},
						"listWithSensitiveSchemaElement": {
							NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
								configschema.DynamicSensitivePathElementKey: {},
							},
						},
						"listWithResourceElement": {},
						"listWithSensitiveResourceElement": {
							NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
								configschema.DynamicSensitivePathElementKey: {
									NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
										"myint": {},
									},
								},
							},
						},
						"listWithValueTypeElement":     {},
						"listWithoutElementDefinition": {},
						"setWithSchemaElement":         {},
						"setWithSensitiveSchemaElement": {
							NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
								configschema.DynamicSensitivePathElementKey: {},
							},
						},
						"setWithResourceElement": {},
						"setWithSensitiveResourceElement": {
							NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
								configschema.DynamicSensitivePathElementKey: {
									NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
										"myint": {},
									},
								},
							},
						},
						"setWithValueTypeElement":     {},
						"setWithoutElementDefinition": {},
						"mapWithSchemaElement":        {},
						"mapWithSensitiveSchemaElement": {
							NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
								configschema.DynamicSensitivePathElementKey: {},
							},
						},
						"unsupported_mapWithResourceElement": {},
						// "unsupported_mapWithSensitiveResourceElement": nil,
						"mapWithValueTypeElement":     {},
						"mapWithoutElementDefinition": {},
					},
				},
			}),
		},
		"sensitive nested attribute": {
			map[string]*Schema{
				"myblock": {
					Type:     TypeList,
					MaxItems: 1,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"attribute": {
								Type:      TypeString,
								Optional:  true,
								Sensitive: true,
							},
						},
					},
				},
			},
			testResource(&configschema.Block{
				Attributes: map[string]*configschema.Attribute{},
				BlockTypes: map[string]*configschema.NestedBlock{
					"myblock": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"attribute": {
									Type:     cty.String,
									Optional: true,
								},
							},
							SensitivePaths: &configschema.SensitivePathElement{
								NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
									"attribute": {},
								},
							},
						},
						Nesting:  configschema.NestingList,
						MaxItems: 1,
					},
				},
				SensitivePaths: &configschema.SensitivePathElement{
					NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
						"myblock": {
							NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
								"attribute": {},
							},
						},
					},
				},
			}),
		},
		"sensitive nested attribute (with explicit ConfigMode)": {
			map[string]*Schema{
				"myblock": {
					Type:       TypeList,
					ConfigMode: SchemaConfigModeBlock,
					MaxItems:   1,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"attribute": {
								Type:      TypeString,
								Optional:  true,
								Sensitive: true,
							},
						},
					},
				},
			},
			testResource(&configschema.Block{
				Attributes: map[string]*configschema.Attribute{},
				BlockTypes: map[string]*configschema.NestedBlock{
					"myblock": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"attribute": {
									Type:     cty.String,
									Optional: true,
								},
							},
							SensitivePaths: &configschema.SensitivePathElement{
								NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
									"attribute": {},
								},
							},
						},
						Nesting:  configschema.NestingList,
						MaxItems: 1,
					},
				},
				SensitivePaths: &configschema.SensitivePathElement{
					NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
						"myblock": {
							NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
								"attribute": {},
							},
						},
					},
				},
			}),
		},
		"required and computed list with nested sensitive value": {
			map[string]*Schema{
				"mylist": {
					Type:     TypeList,
					Computed: true,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"sensitive": {
								Type:      TypeString,
								Sensitive: true,
							},
						},
					},
				},
			},
			testResource(&configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"mylist": {
						Type: cty.List(cty.Object(map[string]cty.Type{
							"sensitive": cty.String,
						})),
						Computed: true,
						Required: true,
					},
				},
				BlockTypes: map[string]*configschema.NestedBlock{},
				SensitivePaths: &configschema.SensitivePathElement{
					NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
						"mylist": {
							NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
								configschema.DynamicSensitivePathElementKey: {
									NestedSensitivePathElements: map[string]*configschema.SensitivePathElement{
										"sensitive": {},
									},
								},
							},
						},
					},
				},
			}),
		},
		//"conditionally required on": {
		//	map[string]*Schema{
		//		"string": {
		//			Type:     TypeString,
		//			Required: true,
		//			DefaultFunc: func() (interface{}, error) {
		//				return nil, nil
		//			},
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{
		//			"string": {
		//				Type:     cty.String,
		//				Required: true,
		//			},
		//		},
		//		BlockTypes: map[string]*configschema.NestedBlock{},
		//	}),
		//},
		//"conditionally required off": {
		//	map[string]*Schema{
		//		"string": {
		//			Type:     TypeString,
		//			Required: true,
		//			DefaultFunc: func() (interface{}, error) {
		//				// If we return a non-nil default then this overrides
		//				// the "Required: true" for the purpose of building
		//				// the core schema, so that core will ignore it not
		//				// being set and let the provider handle it.
		//				return "boop", nil
		//			},
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{
		//			"string": {
		//				Type:     cty.String,
		//				Optional: true,
		//			},
		//		},
		//		BlockTypes: map[string]*configschema.NestedBlock{},
		//	}),
		//},
		//"conditionally required error": {
		//	map[string]*Schema{
		//		"string": {
		//			Type:     TypeString,
		//			Required: true,
		//			DefaultFunc: func() (interface{}, error) {
		//				return nil, fmt.Errorf("placeholder error")
		//			},
		//		},
		//	},
		//	testResource(&configschema.Block{
		//		Attributes: map[string]*configschema.Attribute{
		//			"string": {
		//				Type:     cty.String,
		//				Optional: true, // Just so we can progress to provider-driven validation and return the error there
		//			},
		//		},
		//		BlockTypes: map[string]*configschema.NestedBlock{},
		//	}),
		//},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := (&Resource{Schema: test.Schema}).CoreConfigSchema()
			if !cmp.Equal(got, test.Want, equateEmpty, typeComparer) {
				t.Error(cmp.Diff(test.Want, got, equateEmpty, typeComparer))
			}
		})
	}
}

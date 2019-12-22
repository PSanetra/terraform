package configschema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSensitivePathElement_sensitivePathLeaf_has_no_child_elements(t *testing.T) {
	assert.Equal(t, 0, len(sensitivePathLeaf.NestedSensitivePathElements))
}

func TestSensitivePathElement_NewSensitivePathLeaf_returns_nil_for_non_sensitive_primitive(t *testing.T) {
	var sensitivePath = NewSensitivePathLeaf(false)

	assert.Nil(t, sensitivePath)
}

func TestSensitivePathElement_NewSensitivePathLeaf_returns_non_sensitivePathLeaf_for_sensitive_primitive(t *testing.T) {
	var sensitivePath = NewSensitivePathLeaf(true)

	assert.Same(t, sensitivePathLeaf, sensitivePath)
}

func TestSensitivePathElement_IsSensitive_returns_false_if_path_element_is_nil(t *testing.T) {
	var sensitivePath = NewSensitivePathLeaf(false)

	assert.False(t, sensitivePath.IsSensitive())
}

func TestSensitivePathElement_IsSensitive_returns_false_if_path_element_has_child_elements(t *testing.T) {
	childPath := NewSensitivePathLeaf(true)
	sensitivePath := &SensitivePathElement{
		NestedSensitivePathElements: map[string]*SensitivePathElement{
			"childElement": childPath,
		},
	}

	assert.False(t, sensitivePath.IsSensitive())
}

func TestSensitivePathElement_IsSensitive_returns_true_if_path_element_has_no_child_elements(t *testing.T) {
	sensitivePath := NewSensitivePathLeaf(true)

	assert.True(t, sensitivePath.IsSensitive())
}

func TestSensitivePathElement_ContainsSensitive_returns_false_if_path_element_is_nil(t *testing.T) {
	var sensitivePath *SensitivePathElement = nil

	assert.False(t, sensitivePath.ContainsSensitive())
}

func TestSensitivePathElement_ContainsSensitive_returns_true_if_path_element_has_child_elements(t *testing.T) {
	childPath := NewSensitivePathLeaf(true)
	sensitivePath := &SensitivePathElement{
		NestedSensitivePathElements: map[string]*SensitivePathElement{
			"childElement": childPath,
		},
	}

	assert.True(t, sensitivePath.ContainsSensitive())
}

func TestSensitivePathElement_ContainsSensitive_returns_true_if_path_element_has_no_child_elements(t *testing.T) {
	sensitivePath := NewSensitivePathLeaf(true)

	assert.True(t, sensitivePath.ContainsSensitive())
}

func TestSensitivePathElement_Get_returns_nil_if_path_element_is_nil(t *testing.T) {
	var sensitivePath *SensitivePathElement = nil

	assert.Nil(t, sensitivePath.Get("childElement"))
}

func TestSensitivePathElement_Get_returns_named_child_path_if_current_path_element_is_not_sensitive(t *testing.T) {
	childPath := NewSensitivePathLeaf(true)
	sensitivePath := &SensitivePathElement{
		NestedSensitivePathElements: map[string]*SensitivePathElement{
			"childElement": childPath,
		},
	}

	assert.Same(t, childPath, sensitivePath.Get("childElement"))
	assert.Nil(t, sensitivePath.Get("nonExistingElement"))
}

func TestSensitivePathElement_Get_returns_current_path_element_if_it_is_sensitive(t *testing.T) {
	sensitivePath := NewSensitivePathLeaf(true)

	assert.Same(t, sensitivePath, sensitivePath.Get("randomChildElement"))
}

func TestSensitivePathElement_GetDynamic_returns_nil_if_path_element_is_nil(t *testing.T) {
	var sensitivePath *SensitivePathElement = nil

	assert.Nil(t, sensitivePath.GetDynamic())
}

func TestSensitivePathElement_GetDynamic_returns_child_path_for_empty_key_if_current_path_element_is_not_sensitive(t *testing.T) {
	childPath := NewSensitivePathLeaf(true)
	sensitivePath := &SensitivePathElement{
		NestedSensitivePathElements: map[string]*SensitivePathElement{
			DynamicSensitivePathElementKey: childPath,
		},
	}

	assert.Same(t, childPath, sensitivePath.GetDynamic())
}

func TestSensitivePathElement_GetDynamic_returns_current_path_element_if_it_is_sensitive(t *testing.T) {
	sensitivePath := NewSensitivePathLeaf(true)

	assert.Same(t, sensitivePath, sensitivePath.GetDynamic())
}

func TestSensitivePathElement_Add_returns_sensitive_element_if_current_element_is_sensitive(t *testing.T) {
	sensitivePath := NewSensitivePathLeaf(true)

	sensitivePath = sensitivePath.Add("nonSensitive", NewSensitivePathLeaf(false))

	assert.True(t, sensitivePath.IsSensitive())
	assert.True(t, sensitivePath.Get("nonSensitive").IsSensitive())
}

func TestSensitivePathElement_Add_returns_element_which_does_not_contain_sensitive_paths_if_non_sensitive_is_added(t *testing.T) {
	sensitivePath := NewSensitivePathLeaf(false)

	sensitivePath = sensitivePath.Add("nonSensitive", NewSensitivePathLeaf(false))

	assert.False(t, sensitivePath.ContainsSensitive())
}

func TestSensitivePathElement_Add_returns_element_which_contains_sensitive_paths_if_sensitive_is_added(t *testing.T) {
	sensitivePath := NewSensitivePathLeaf(false)

	sensitivePath = sensitivePath.Add("sensitive", NewSensitivePathLeaf(true))

	assert.False(t, sensitivePath.IsSensitive())
	assert.True(t, sensitivePath.ContainsSensitive())
	assert.True(t, sensitivePath.Get("sensitive").IsSensitive())
}

func TestSensitivePathElement_Add_returns_element_which_contains_sensitive_paths_if_element_containing_sensitive_is_added(t *testing.T) {
	sensitivePath := NewSensitivePathLeaf(false)

	sensitivePath = sensitivePath.Add(
		"containsSensitive",
		NewSensitivePathLeaf(false).Add(
			"sensitive",
			NewSensitivePathLeaf(true),
		),
	)

	assert.False(t, sensitivePath.IsSensitive())
	assert.True(t, sensitivePath.ContainsSensitive())
	assert.False(t, sensitivePath.Get("containsSensitive").IsSensitive())
	assert.True(t, sensitivePath.Get("containsSensitive").ContainsSensitive())
	assert.True(t, sensitivePath.Get("containsSensitive").Get("sensitive").IsSensitive())
}

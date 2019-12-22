package configschema

const DynamicSensitivePathElementKey string = ""

// Contains paths to sensitive values.
// Only the leafs of this path element tree are considered sensitive
// as all possible children of a sensitive element would be sensitive, too,
type SensitivePathElement struct {
	// The key of this map is an empty string if the child elements have dynamic names like map keys or list indices
	// The key is non empty for child elements with static names like attribute names
	NestedSensitivePathElements map[string]*SensitivePathElement
}

var nonSensitivePath *SensitivePathElement = nil
var sensitivePathLeaf = &SensitivePathElement{NestedSensitivePathElements: nil}

func NewSensitivePathLeaf(forSensitivePrimitive bool) *SensitivePathElement {
	if forSensitivePrimitive {
		return sensitivePathLeaf
	} else {
		return nonSensitivePath
	}
}

func (s *SensitivePathElement) Add(childElement string, paths *SensitivePathElement) *SensitivePathElement {
	if s.IsSensitive() || !paths.ContainsSensitive() {
		return s
	}

	if s == nil {
		return &SensitivePathElement{
			NestedSensitivePathElements: map[string]*SensitivePathElement{
				childElement: paths,
			},
		}
	}

	s.NestedSensitivePathElements[childElement] = paths

	return s
}

func (s *SensitivePathElement) AddDynamic(paths *SensitivePathElement) *SensitivePathElement {
	return s.Add(DynamicSensitivePathElementKey, paths)
}

func (s *SensitivePathElement) Get(childElement string) *SensitivePathElement {
	if s == nil {
		return nil
	}

	if s.IsSensitive() {
		return s
	}

	return s.NestedSensitivePathElements[childElement]
}

// Returns a sensitive path element which represents all dynamically named child paths (e.g. list or map elements)
func (s *SensitivePathElement) GetDynamic() *SensitivePathElement {
	return s.Get(DynamicSensitivePathElementKey)
}

// Returns true if this path element is considered sensitive
func (s *SensitivePathElement) IsSensitive() bool {
	if s == nil {
		return false
	}

	// An element of a sensitive path is considered sensitive if it is the leaf path element
	return len(s.NestedSensitivePathElements) == 0
}

// Returns true if this element or a descendant element is sensitive
func (s *SensitivePathElement) ContainsSensitive() bool {
	return s != nil
}

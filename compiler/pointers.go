package compiler

// StrPtr returns the pointer of a string.
func StrPtr(v string) *string {
	return &v
}

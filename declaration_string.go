// Code generated by "stringer -type Declaration -trimprefix Decl"; DO NOT EDIT.

package gosort

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DeclOther-0]
	_ = x[DeclConstructor-1]
	_ = x[DeclType-2]
	_ = x[DeclMethod-3]
	_ = x[DeclFunc-4]
}

const _Declaration_name = "OtherConstructorTypeMethodFunc"

var _Declaration_index = [...]uint8{0, 5, 16, 20, 26, 30}

func (i Declaration) String() string {
	if i < 0 || i >= Declaration(len(_Declaration_index)-1) {
		return "Declaration(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Declaration_name[_Declaration_index[i]:_Declaration_index[i+1]]
}

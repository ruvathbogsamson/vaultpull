package audit

import (
	"path/filepath"
	"sort"
)

// This file provides the imports used by rotation.go's pruneBackups method.
// filepath and sort are referenced there; Go requires them in the same package.

var (
	_ = filepath.Glob
	_ = sort.Strings
)

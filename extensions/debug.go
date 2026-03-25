package extensions

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func toJSON(v any) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}

func DumpArgs(args ...any) (string, error) {
	var result strings.Builder
	if args != nil {
		for i, arg := range args {
			result.WriteString(fmt.Sprintf("%d => %s (%T)\n", i, toJSON(arg), arg))
		}
		fmt.Fprintln(os.Stderr, result.String())
		return result.String(), nil
	} else {
		return "<empty>", nil
	}
}

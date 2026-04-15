package preview

import (
	"regexp"
	"sort"
	"strings"
)

var placeholderRE = regexp.MustCompile(`\{\{\s*([^}]+?)\s*\}\}`)

func Resolve(s string, vars map[string]string) (string, []string) {
	seenMissing := map[string]struct{}{}
	var missing []string

	out := placeholderRE.ReplaceAllStringFunc(s, func(full string) string {
		sub := placeholderRE.FindStringSubmatch(full)
		if len(sub) < 2 {
			return full
		}
		key := strings.TrimSpace(sub[1])
		if val, ok := vars[key]; ok {
			return val
		}
		if _, ok := seenMissing[key]; !ok {
			seenMissing[key] = struct{}{}
			missing = append(missing, key)
		}
		return full
	})

	sort.Strings(missing)
	return out, missing
}

func MergeMissing(parts ...[]string) []string {
	m := map[string]struct{}{}
	for _, list := range parts {
		for _, k := range list {
			m[k] = struct{}{}
		}
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

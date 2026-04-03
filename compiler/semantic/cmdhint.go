package semantic

import (
	"fmt"
	"strings"

	"moonbasic/compiler/builtinmanifest"
)

// unknownCommandMessageAndHint builds user-facing text for a bad NS.METHOD reference.
func unknownCommandMessageAndHint(t *builtinmanifest.Table, ns, method string) (msg, hint string) {
	full := builtinmanifest.Key(ns, method)
	msg = fmt.Sprintf("Unknown command '%s'", full)

	inNS := t.KeysWithNamespacePrefix(ns)
	if len(inNS) > 0 {
		bestM := ""
		bestD := 99
		suffix := ns + "."
		for _, k := range inNS {
			m := strings.TrimPrefix(k, suffix)
			d := builtinmanifest.EditDistance(method, m)
			if d < bestD {
				bestD = d
				bestM = k
			}
		}
		if bestD <= 3 && bestM != "" {
			hint = fmt.Sprintf("Did you mean %s ?", bestM)
			list := builtinmanifest.FormatNamespaceListing(ns, inNS, 76)
			if list != "" {
				hint += "\n" + list
			}
			return msg, hint
		}
		list := builtinmanifest.FormatNamespaceListing(ns, inNS, 76)
		if list != "" {
			return msg, list
		}
	}

	if alt, ok := t.BestSimilarKey(full, 4); ok {
		return msg, fmt.Sprintf("Did you mean %s ?", alt)
	}

	return msg, "Add the command to compiler/builtinmanifest/commands.json or use a supported API."
}

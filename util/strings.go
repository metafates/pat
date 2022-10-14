package util

import (
	"fmt"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"strings"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}

func Quantify(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}

	return fmt.Sprintf("%d %s", n, plural)
}

func FindClosest(s string, ss []string) mo.Option[string] {
	closest := lo.MaxBy(ss, func(a, b string) bool {
		return levenshtein.Distance(a, s) < levenshtein.Distance(b, s)
	})

	if levenshtein.Distance(closest, s) > 3 {
		return mo.None[string]()
	}

	return mo.Some(closest)
}

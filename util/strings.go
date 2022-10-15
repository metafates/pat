package util

import (
	"fmt"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"os"
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

func CompareSemVers(a, b string) (int, error) {
	type version struct {
		major, minor, patch int
	}

	parse := func(s string) (version, error) {
		var v version
		_, err := fmt.Sscanf(strings.TrimPrefix(s, "v"), "%d.%d.%d", &v.major, &v.minor, &v.patch)
		return v, err
	}

	av, err := parse(a)
	if err != nil {
		return 0, err
	}

	bv, err := parse(b)
	if err != nil {
		return 0, err
	}

	for _, pair := range []lo.Tuple2[int, int]{
		{av.major, bv.major},
		{av.minor, bv.minor},
		{av.patch, bv.patch},
	} {
		if pair.A > pair.B {
			return 1, nil
		}

		if pair.A < pair.B {
			return -1, nil
		}
	}

	return 0, nil
}

// PrintErasable prints a string that can be erased by calling a returned function.
func PrintErasable(msg string) (eraser func()) {
	_, _ = fmt.Fprintf(os.Stdout, "\r%s", msg)

	return func() {
		_, _ = fmt.Fprintf(os.Stdout, "\r%s\r", strings.Repeat(" ", len(msg)))
	}
}

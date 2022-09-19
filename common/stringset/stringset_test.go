package stringset

import (
	"testing"
)

func TestStringSet_Match(t *testing.T) {

	set := New()
	set.Add("aaaaa", "bbbbbb", "ccc", "ddd")
	t.Log(set.Contains("aaaaa"))
	t.Log(set.Contains("ccc"))
	t.Log(set.MatchFirst("aasdddccasdadsads"))
	set.Remove("ddd")
	t.Log(set.MatchFirst("aasdddccasdadsads"))
}

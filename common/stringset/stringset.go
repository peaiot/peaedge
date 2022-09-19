package stringset

import (
	"fmt"
	"strings"
	"sync"
)

// Set is an unsorted set of unique strings.
type StringSet struct {
	sync.RWMutex
	Set map[string]struct{}
}

// New returns an initialized Set.
func New() *StringSet {
	return &StringSet{Set: make(map[string]struct{})}
}

// NewFromSlice returns a new set constructed from the given slice. Any
// duplicate elements will be removed.
func NewFromSlice(slice []string) *StringSet {
	s := New()
	for _, v := range slice {
		s.Add(v)
	}
	return s
}

func (s *StringSet) Copy() *StringSet {
	ns := New()
	s.Lock()
	defer s.Unlock()
	for v := range s.Set {
		ns.Set[v] = struct{}{}
	}
	return ns
}

// Add adds each value in vs to the set. If any value is alredy in the set,
// this has no effect.
func (s *StringSet) Add(vs ...string) {
	s.Lock()
	defer s.Unlock()
	for _, v := range vs {
		s.Set[v] = struct{}{}
	}
}

// Remove removes v from the set. If v is not in the set, this has no effect.
func (s *StringSet) Remove(v string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Set, v)
}

// Contains returns true if the set contains v and false otherwise.
func (s *StringSet) Contains(v string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.Set[v]
	return ok
}

func (s *StringSet) MatchFirst(src string) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	for v := range s.Set {
		if strings.Contains(src, v) {
			return v, true
		}
	}
	return "", false
}

// Slice returns the elements in the set as a slice of strings. It returns an
// empty slice if the set contains no elements. The elements returned will be
// in random order.
func (s *StringSet) Slice() []string {
	s.RLock()
	defer s.RUnlock()
	slice := make([]string, len(s.Set))
	i := 0
	for v := range s.Set {
		slice[i] = v
		i++
	}
	return slice
}

// String implements the Stringer interface.
func (s *StringSet) String() string {
	return fmt.Sprint(s.Slice())
}

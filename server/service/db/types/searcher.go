package types

type Searcher interface {
	Conditions() []Condition
}

type searcherOption interface {
	apply(*searcher)
}

type searcherOptionFunc func(*searcher)

func (fn searcherOptionFunc) apply(s *searcher) {
	fn(s)
}

func WithSearchMap(searchMap map[string]string) searcherOption {
	return searcherOptionFunc(func(s *searcher) {
		s.searchMap = searchMap
	})
}

func WithValidColumnsForSearcher(validColumns map[string]struct{}) searcherOption {
	return searcherOptionFunc(func(s *searcher) {
		s.validColumns = validColumns
	})
}

type searcher struct {
	conditions   []Condition
	searchMap    map[string]string
	validColumns map[string]struct{}
}

func (s *searcher) Conditions() []Condition {
	if s.conditions == nil {
		s.conditions = make([]Condition, 0)

		for column, value := range s.searchMap {
			if _, ok := s.validColumns[column]; !ok {
				continue
			}

			s.conditions = append(s.conditions, OrCondition(column, value))
		}
	}

	return s.conditions
}

func NewSearcher(searcherOpts ...searcherOption) *searcher {
	s := &searcher{
		searchMap:    make(map[string]string),
		validColumns: make(map[string]struct{}),
	}

	for _, opt := range searcherOpts {
		opt.apply(s)
	}

	return s
}

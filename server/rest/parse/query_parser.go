package parse

import (
	"net/url"
	"strconv"
	"strings"

	"space.online.shop.web.server/service/db/types"
)

const (
	PAGE      string = "page"
	PAGE_SIZE string = "page_size"
)

type queryParser struct {
	queryMap     map[string]string
	counter      types.Counter
	sortOrder    types.SortOrder
	searcher     types.Searcher
	validColumns map[string]struct{}
}

func (q *queryParser) Counter() types.Counter {
	if q.counter == nil {
		page, pageSize := q.parseCounter()

		q.counter = types.NewCounter(
			types.WithPage(page),
			types.WithPageSize(pageSize),
		)
	}
	return q.counter
}

func (q *queryParser) parseCounter() (int, int) {
	page := types.DEFAULT_PAGE
	if pageStr, ok := q.queryMap[PAGE]; ok {
		parsedPage, err := strconv.ParseInt(pageStr, 10, 64)
		if err == nil {
			page = int(parsedPage)
		}
	}

	pageSize := types.DEFAULT_PAGE_SIZE
	if pageSizeStr, ok := q.queryMap[PAGE_SIZE]; ok {
		parsedPageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err == nil {
			pageSize = int(parsedPageSize)
		}
	}

	return page, pageSize
}

func (q *queryParser) SortOrder() types.SortOrder {
	if q.sortOrder == nil {
		sortOrderMap := q.parseSortOrderMap()

		q.sortOrder = types.NewSortOrder(
			types.WithSortOrderMap(sortOrderMap),
			types.WithValidColumnsForSortOrder(q.validColumns),
		)
	}

	return q.sortOrder
}

func (q *queryParser) parseSortOrderMap() map[string]string {
	sortOrderMap := make(map[string]string)

	for k, v := range q.queryMap {
		if strings.HasSuffix(k, "_sort_order") {
			column := strings.TrimSuffix(k, "_sort_order")
			if column == "" {
				continue
			}
			sortOrderMap[column] = v
		}
	}

	return sortOrderMap
}

func (q *queryParser) Searcher() types.Searcher {
	if q.searcher == nil {
		searchMap := q.parseSearchMap()

		q.searcher = types.NewSearcher(
			types.WithSearchMap(searchMap),
			types.WithValidColumnsForSearcher(q.validColumns),
		)
	}

	return q.searcher
}

func (q *queryParser) parseSearchMap() map[string]string {
	searchMap := make(map[string]string)

	for k, v := range q.queryMap {
		if strings.HasSuffix(k, "_sort_order") {
			continue
		}

		switch k {
		case PAGE:
			continue
		case PAGE_SIZE:
			continue
		}

		searchMap[k] = v
	}

	return searchMap
}

func NewQueryParser(value url.Values, validColumns map[string]struct{}) *queryParser {
	parser := &queryParser{
		queryMap:     make(map[string]string),
		validColumns: validColumns,
	}

	for k, v := range value {
		if len(v) < 1 {
			parser.queryMap[k] = ""
		} else {
			parser.queryMap[k] = v[0]
		}
	}

	return parser
}

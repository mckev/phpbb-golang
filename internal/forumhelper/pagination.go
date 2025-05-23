package forumhelper

import (
	"maps"
	"slices"
)

// PaginationType has 5 types:
//   - Arrow Previous: start item
//   - Page: start item, page number
//   - Current Page: page number
//   - Separator
//   - Arrow Next: start item
const (
	PAGINATION_TYPE_ARROW_PREVIOUS = "PaginationTypeArrowPrevious"
	PAGINATION_TYPE_PAGE           = "PaginationTypePage"
	PAGINATION_TYPE_CURRENT_PAGE   = "PaginationTypeCurrentPage"
	PAGINATION_TYPE_SEPARATOR      = "PaginationTypeSeparator"
	PAGINATION_TYPE_ARROW_NEXT     = "PaginationTypeArrowNext"
)

type Pagination struct {
	PaginationType string `json:"pagination_type"`
	StartItem      int    `json:"start_item"`  // StartItem strats at 0
	PageNumber     int    `json:"page_number"` // PageNumber starts at 1
}

func ComputePaginations(curItem int, totalItems int, maxItemsPerPage int) []Pagination {
	// Notes:
	//   - curItem starts at 0, and pages starts at 0
	//   - curPage can be equal to maxPage
	curPage := curItem / maxItemsPerPage
	maxPage := totalItems / maxItemsPerPage

	// Process pages around curPage, which we call "midPages"
	midPagesSet := map[int]bool{}
	for page := 0; page <= maxPage; page++ { // Equivalent:	for page := range maxPage + 1 {
		if curPage+page <= maxPage {
			midPagesSet[curPage+page] = true
		}
		if len(midPagesSet) >= 5 {
			break
		}
		if curPage-page >= 0 {
			midPagesSet[curPage-page] = true
		}
		if len(midPagesSet) >= 5 {
			break
		}
	}
	midPagesList := slices.Collect(maps.Keys(midPagesSet))
	slices.Sort(midPagesList)
	minMidPages := midPagesList[0]
	maxMidPages := midPagesList[len(midPagesList)-1]

	// Construct Paginations
	paginations := []Pagination{}
	if curPage-1 >= 0 {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_ARROW_PREVIOUS,
			StartItem:      (curPage - 1) * maxItemsPerPage,
		})
	}
	if minMidPages >= 1 {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_PAGE,
			StartItem:      0,
			PageNumber:     1,
		})
	}
	if minMidPages >= 2 {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_SEPARATOR,
		})
	}
	for _, page := range midPagesList {
		if page == curPage {
			paginations = append(paginations, Pagination{
				PaginationType: PAGINATION_TYPE_CURRENT_PAGE,
				PageNumber:     page + 1,
			})
		} else {
			paginations = append(paginations, Pagination{
				PaginationType: PAGINATION_TYPE_PAGE,
				StartItem:      page * maxItemsPerPage,
				PageNumber:     page + 1,
			})
		}
	}
	if maxMidPages <= maxPage-2 {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_SEPARATOR,
		})
	}
	if maxMidPages <= maxPage-1 {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_PAGE,
			StartItem:      maxPage * maxItemsPerPage,
			PageNumber:     maxPage + 1,
		})
	}
	if curPage+1 <= maxPage {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_ARROW_NEXT,
			StartItem:      (curPage + 1) * maxItemsPerPage,
		})
	}
	return paginations
}

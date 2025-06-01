package forumhelper

import (
	"maps"
	"slices"
)

// PaginationType has 5 types:
//   - Arrow Previous: contains start item
//   - Page: contains start item, page number
//   - Current Page: contains start item, page number
//   - Separator
//   - Arrow Next: contains start item
const (
	PAGINATION_TYPE_ARROW_PREVIOUS = "PaginationTypeArrowPrevious"
	PAGINATION_TYPE_PAGE           = "PaginationTypePage"
	PAGINATION_TYPE_CURRENT_PAGE   = "PaginationTypeCurrentPage"
	PAGINATION_TYPE_SEPARATOR      = "PaginationTypeSeparator"
	PAGINATION_TYPE_ARROW_NEXT     = "PaginationTypeArrowNext"
)

type Pagination struct {
	PaginationType string `json:"pagination_type"`
	StartItem      int    `json:"start_item"`  // StartItem starts at 0
	PageNumber     int    `json:"page_number"` // PageNumber starts at 1
}

func ComputePaginations(curItem int, totalItems int, maxItemsPerPage int) []Pagination {
	// Notes:
	//   - 0 <= curItem < totalItems
	//   - 0 <= page < maxPage
	//   - page starts at 0, however page number (i.e. human-readable page) starts at 1
	if totalItems == 0 {
		return []Pagination{}
	}
	curPage := curItem / maxItemsPerPage
	maxPage := (totalItems + maxItemsPerPage - 1) / maxItemsPerPage // maxPage := ceil(totalItems / maxItemsPerPage)

	// Process pages around curPage, which we call "midPages"
	midPagesSet := map[int]bool{}
	for page := 0; page < maxPage; page++ { // Modern version:  for page := range maxPage {
		if curPage+page < maxPage {
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
				StartItem:      page * maxItemsPerPage,
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
	if maxMidPages <= maxPage-3 {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_SEPARATOR,
		})
	}
	if maxMidPages <= maxPage-2 {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_PAGE,
			StartItem:      (maxPage - 1) * maxItemsPerPage,
			PageNumber:     maxPage,
		})
	}
	if curPage+1 <= maxPage-1 {
		paginations = append(paginations, Pagination{
			PaginationType: PAGINATION_TYPE_ARROW_NEXT,
			StartItem:      (curPage + 1) * maxItemsPerPage,
		})
	}
	return paginations
}

func ComputeStartItem(curItem int, totalItems int, maxItemsPerPage int) int {
	if curItem < 0 {
		curItem = 0
	}
	if curItem > totalItems-1 {
		curItem = totalItems - 1
	}
	startItem := (curItem / maxItemsPerPage) * maxItemsPerPage
	return startItem
}

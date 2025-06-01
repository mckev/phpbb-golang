package forumhelper

import (
	"reflect"
	"testing"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/model"
)

func TestComputePaginations_NoElement(t *testing.T) {
	actual := forumhelper.ComputePaginations(0, 0, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_SingleElement(t *testing.T) {
	actual := forumhelper.ComputePaginations(0, 1, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 0, PageNumber: 1},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page1of24(t *testing.T) {
	actual := forumhelper.ComputePaginations(5, 580, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		// We should not have an Arrow Previous here
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 0, PageNumber: 1},
		{PaginationType: "PaginationTypePage", StartItem: 25, PageNumber: 2},
		{PaginationType: "PaginationTypePage", StartItem: 50, PageNumber: 3},
		{PaginationType: "PaginationTypePage", StartItem: 75, PageNumber: 4},
		{PaginationType: "PaginationTypePage", StartItem: 100, PageNumber: 5},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 575, PageNumber: 24},
		{PaginationType: "PaginationTypeArrowNext", StartItem: 25, PageNumber: 0},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page2of24(t *testing.T) {
	actual := forumhelper.ComputePaginations(30, 580, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeArrowPrevious", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 0, PageNumber: 1},
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 25, PageNumber: 2},
		{PaginationType: "PaginationTypePage", StartItem: 50, PageNumber: 3},
		{PaginationType: "PaginationTypePage", StartItem: 75, PageNumber: 4},
		{PaginationType: "PaginationTypePage", StartItem: 100, PageNumber: 5},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 575, PageNumber: 24},
		{PaginationType: "PaginationTypeArrowNext", StartItem: 50, PageNumber: 0},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page4of24(t *testing.T) {
	actual := forumhelper.ComputePaginations(80, 580, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeArrowPrevious", StartItem: 50, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 0, PageNumber: 1},
		// We should not have a Separator here
		{PaginationType: "PaginationTypePage", StartItem: 25, PageNumber: 2},
		{PaginationType: "PaginationTypePage", StartItem: 50, PageNumber: 3},
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 75, PageNumber: 4},
		{PaginationType: "PaginationTypePage", StartItem: 100, PageNumber: 5},
		{PaginationType: "PaginationTypePage", StartItem: 125, PageNumber: 6},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 575, PageNumber: 24},
		{PaginationType: "PaginationTypeArrowNext", StartItem: 100, PageNumber: 0},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page5of24(t *testing.T) {
	actual := forumhelper.ComputePaginations(105, 580, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeArrowPrevious", StartItem: 75, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 0, PageNumber: 1},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 50, PageNumber: 3},
		{PaginationType: "PaginationTypePage", StartItem: 75, PageNumber: 4},
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 100, PageNumber: 5},
		{PaginationType: "PaginationTypePage", StartItem: 125, PageNumber: 6},
		{PaginationType: "PaginationTypePage", StartItem: 150, PageNumber: 7},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 575, PageNumber: 24},
		{PaginationType: "PaginationTypeArrowNext", StartItem: 125, PageNumber: 0},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page18of24(t *testing.T) {
	actual := forumhelper.ComputePaginations(430, 580, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeArrowPrevious", StartItem: 400, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 0, PageNumber: 1},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 375, PageNumber: 16},
		{PaginationType: "PaginationTypePage", StartItem: 400, PageNumber: 17},
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 425, PageNumber: 18},
		{PaginationType: "PaginationTypePage", StartItem: 450, PageNumber: 19},
		{PaginationType: "PaginationTypePage", StartItem: 475, PageNumber: 20},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 575, PageNumber: 24},
		{PaginationType: "PaginationTypeArrowNext", StartItem: 450, PageNumber: 0},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page21of24(t *testing.T) {
	actual := forumhelper.ComputePaginations(505, 580, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeArrowPrevious", StartItem: 475, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 0, PageNumber: 1},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 450, PageNumber: 19},
		{PaginationType: "PaginationTypePage", StartItem: 475, PageNumber: 20},
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 500, PageNumber: 21},
		{PaginationType: "PaginationTypePage", StartItem: 525, PageNumber: 22},
		{PaginationType: "PaginationTypePage", StartItem: 550, PageNumber: 23},
		// We should not have a Separator here
		{PaginationType: "PaginationTypePage", StartItem: 575, PageNumber: 24},
		{PaginationType: "PaginationTypeArrowNext", StartItem: 525, PageNumber: 0},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page23of24(t *testing.T) {
	actual := forumhelper.ComputePaginations(555, 580, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeArrowPrevious", StartItem: 525, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 0, PageNumber: 1},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 475, PageNumber: 20},
		{PaginationType: "PaginationTypePage", StartItem: 500, PageNumber: 21},
		{PaginationType: "PaginationTypePage", StartItem: 525, PageNumber: 22},
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 550, PageNumber: 23},
		{PaginationType: "PaginationTypePage", StartItem: 575, PageNumber: 24},
		{PaginationType: "PaginationTypeArrowNext", StartItem: 575, PageNumber: 0},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page24of24(t *testing.T) {
	actual := forumhelper.ComputePaginations(576, 580, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeArrowPrevious", StartItem: 550, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 0, PageNumber: 1},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 475, PageNumber: 20},
		{PaginationType: "PaginationTypePage", StartItem: 500, PageNumber: 21},
		{PaginationType: "PaginationTypePage", StartItem: 525, PageNumber: 22},
		{PaginationType: "PaginationTypePage", StartItem: 550, PageNumber: 23},
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 575, PageNumber: 24},
		// We should not have an Arrow Next here
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputePaginations_Page8of8(t *testing.T) {
	actual := forumhelper.ComputePaginations(180, 200, model.MAX_POSTS_PER_PAGE)
	expected := []forumhelper.Pagination{
		{PaginationType: "PaginationTypeArrowPrevious", StartItem: 150, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 0, PageNumber: 1},
		{PaginationType: "PaginationTypeSeparator", StartItem: 0, PageNumber: 0},
		{PaginationType: "PaginationTypePage", StartItem: 75, PageNumber: 4},
		{PaginationType: "PaginationTypePage", StartItem: 100, PageNumber: 5},
		{PaginationType: "PaginationTypePage", StartItem: 125, PageNumber: 6},
		{PaginationType: "PaginationTypePage", StartItem: 150, PageNumber: 7},
		{PaginationType: "PaginationTypeCurrentPage", StartItem: 175, PageNumber: 8},
		// We should not have an Arrow Next here
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, wanted %v", actual, expected)
		return
	}
}

func TestComputeStartItem(t *testing.T) {
	{
		actual := forumhelper.ComputeStartItem(-1000, 250, 25)
		expected := 0
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeStartItem(0, 250, 25)
		expected := 0
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeStartItem(24, 250, 25)
		expected := 0
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeStartItem(25, 250, 25)
		expected := 25
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeStartItem(26, 250, 25)
		expected := 25
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeStartItem(249, 250, 25)
		expected := 225
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeStartItem(250, 250, 25)
		expected := 225
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeStartItem(251, 250, 25)
		expected := 225
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeStartItem(2000, 250, 25)
		expected := 225
		if actual != expected {
			t.Errorf("Got %d, wanted %d", actual, expected)
			return
		}
	}
}

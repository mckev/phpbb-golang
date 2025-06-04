package forumhelper

import (
	"context"
	"reflect"
	"testing"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/model"
)

func TestComputeForumNavTrails(t *testing.T) {
	ctx := context.Background()
	{
		actual := forumhelper.ComputeForumNavTrails(ctx, forums, model.ROOT_FORUM_ID)
		expected := []forumhelper.ForumNavTrail{}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeForumNavTrails(ctx, forums, 1)
		expected := []forumhelper.ForumNavTrail{
			{Forum: model.Forum{ForumId: 1, ParentId: 0, ForumName: "Forum A"}, IsLeaf: true},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeForumNavTrails(ctx, forums, 6)
		expected := []forumhelper.ForumNavTrail{
			{Forum: model.Forum{ForumId: 2, ParentId: 0, ForumName: "Forum B"}, IsLeaf: false},
			{Forum: model.Forum{ForumId: 4, ParentId: 2, ForumName: "Forum B1"}, IsLeaf: false},
			{Forum: model.Forum{ForumId: 6, ParentId: 4, ForumName: "Forum B1A"}, IsLeaf: true},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
}

func TestComputeForumNavTrails_Invalid(t *testing.T) {
	ctx := context.Background()
	{
		actual := forumhelper.ComputeForumNavTrails(ctx, forums, model.INVALID_FORUM_ID)
		expected := []forumhelper.ForumNavTrail{}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
}

func TestComputeForumChildNodes(t *testing.T) {
	ctx := context.Background()
	{
		actual := forumhelper.ComputeForumChildNodes(ctx, forums, model.ROOT_FORUM_ID, 0)
		expected := []forumhelper.ForumNode{
			{Forum: model.Forum{ForumId: 1, ParentId: 0, ForumName: "Forum A"}, IsLeaf: true, ForumChildNodes: []forumhelper.ForumNode{}},
			{Forum: model.Forum{ForumId: 2, ParentId: 0, ForumName: "Forum B"}, IsLeaf: false, ForumChildNodes: []forumhelper.ForumNode{
				{Forum: model.Forum{ForumId: 4, ParentId: 2, ForumName: "Forum B1"}, IsLeaf: false, ForumChildNodes: []forumhelper.ForumNode{
					{Forum: model.Forum{ForumId: 6, ParentId: 4, ForumName: "Forum B1A"}, IsLeaf: true, ForumChildNodes: []forumhelper.ForumNode{}},
				}},
				{Forum: model.Forum{ForumId: 5, ParentId: 2, ForumName: "Forum B2"}, IsLeaf: true, ForumChildNodes: []forumhelper.ForumNode{}},
			}},
			{Forum: model.Forum{ForumId: 3, ParentId: 0, ForumName: "Forum C"}, IsLeaf: true, ForumChildNodes: []forumhelper.ForumNode{}},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeForumChildNodes(ctx, forums, 1, 0)
		expected := []forumhelper.ForumNode{}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeForumChildNodes(ctx, forums, 2, 0)
		expected := []forumhelper.ForumNode{
			{Forum: model.Forum{ForumId: 4, ParentId: 2, ForumName: "Forum B1"}, IsLeaf: false, ForumChildNodes: []forumhelper.ForumNode{
				{Forum: model.Forum{ForumId: 6, ParentId: 4, ForumName: "Forum B1A"}, IsLeaf: true, ForumChildNodes: []forumhelper.ForumNode{}},
			}},
			{Forum: model.Forum{ForumId: 5, ParentId: 2, ForumName: "Forum B2"}, IsLeaf: true, ForumChildNodes: []forumhelper.ForumNode{}},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
	{
		actual := forumhelper.ComputeForumChildNodes(ctx, forums, 4, 0)
		expected := []forumhelper.ForumNode{
			{Forum: model.Forum{ForumId: 6, ParentId: 4, ForumName: "Forum B1A"}, IsLeaf: true, ForumChildNodes: []forumhelper.ForumNode{}},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
}

func TestComputeForumChildNodes_Invalid(t *testing.T) {
	ctx := context.Background()
	{
		actual := forumhelper.ComputeForumChildNodes(ctx, forums, model.INVALID_FORUM_ID, 0)
		expected := []forumhelper.ForumNode{}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Got %v, wanted %v", actual, expected)
			return
		}
	}
}

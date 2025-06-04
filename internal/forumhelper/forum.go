package forumhelper

import (
	"context"
	"slices"

	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

type ForumNavTrail struct {
	Forum  model.Forum
	IsLeaf bool
}

func ComputeForumNavTrails(ctx context.Context, forums []model.Forum, forumId int) []ForumNavTrail {
	// Given a Forum Id, find its parents until Root Forum
	// Convert forums from a list into a map
	forumsMap := map[int]model.Forum{}
	for _, forum := range forums {
		forumsMap[forum.ForumId] = forum
	}
	forumNavTrails := []ForumNavTrail{}
	depth := 0
	id := forumId
	for true {
		if id == model.ROOT_FORUM_ID {
			break
		}
		forum, ok := forumsMap[id]
		if !ok {
			logger.Warnf(ctx, "Error while computing Forum Nav Trails for forum id %d: Unknown forum id %d", forumId, id)
			return []ForumNavTrail{}
		}
		forumNavTrails = append(forumNavTrails, ForumNavTrail{
			Forum:  forum,
			IsLeaf: false,
		})
		depth++
		if depth > model.MAX_FORUM_NAV_TRAILS_DEPTH {
			logger.Warnf(ctx, "Error while computing Forum Nav Trails for forum id %d: Path too deep", forumId)
			return []ForumNavTrail{}
		}
		id = forum.ParentId
	}
	if len(forumNavTrails) > 0 {
		slices.Reverse(forumNavTrails)
		forumNavTrails[len(forumNavTrails)-1].IsLeaf = true
	}
	return forumNavTrails
}

type ForumNode struct {
	Forum           model.Forum
	IsLeaf          bool
	ForumChildNodes []ForumNode
}

func ComputeForumChildNodes(ctx context.Context, forums []model.Forum, forumId int, depth int) []ForumNode {
	// Construct child nodes whose parent is Forum Id
	if depth >= model.MAX_FORUM_NODES_DEPTH {
		return []ForumNode{}
	}
	forumChildNodes := []ForumNode{}
	for _, forum := range forums {
		if forum.ForumId == model.ROOT_FORUM_ID {
			continue
		}
		if forum.ParentId == forumId {
			forumGrandchildNodes := ComputeForumChildNodes(ctx, forums, forum.ForumId, depth+1)
			isLeaf := false
			if len(forumGrandchildNodes) == 0 {
				isLeaf = true
			}
			forumChildNodes = append(forumChildNodes, ForumNode{
				Forum:           forum,
				IsLeaf:          isLeaf,
				ForumChildNodes: forumGrandchildNodes,
			})
		}
	}
	return forumChildNodes
}

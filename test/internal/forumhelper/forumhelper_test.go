package forumhelper

import (
	"os"
	"testing"

	"phpbb-golang/model"
)

var forums []model.Forum

func TestMain(m *testing.M) {
	// Set up
	forums = []model.Forum{
		{ForumId: 0, ParentId: 0, ForumName: "Root Forum"},
		{ForumId: 1, ParentId: 0, ForumName: "Forum A"},
		{ForumId: 2, ParentId: 0, ForumName: "Forum B"},
		{ForumId: 3, ParentId: 0, ForumName: "Forum C"},
		{ForumId: 4, ParentId: 2, ForumName: "Forum B1"},
		{ForumId: 5, ParentId: 2, ForumName: "Forum B2"},
		{ForumId: 6, ParentId: 4, ForumName: "Forum B1A"},
	}
	// Run tests
	exitVal := m.Run()
	// Tear down
	// Exit
	os.Exit(exitVal)
}

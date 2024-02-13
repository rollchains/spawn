package spawn_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	spawn "gitub.com/strangelove-ventures/spawn/spawn"
)

func TestUncomment(t *testing.T) {
	content := `SourceCode goes here!
//this is line 1 which should be uncommented // ?spawntag::test
untouched line 3 // spawntag:test`

	fc := &spawn.FileContent{
		Contents: content,
	}

	fc.HandleCommentSwaps("test")

	// the line should be uncommented
	require.Equal(t, "this is line 1 which should be uncommented", strings.Split(fc.Contents, "\n")[1])
}

func TestRemoveLine(t *testing.T) {
	content := `SourceCode goes here!
//this is line 1 which should be uncommented // ?spawntag::test
this line gets deleted //spawntag:test`

	fc := &spawn.FileContent{
		Contents: content,
	}

	deleteLine := true
	fc.RemoveTaggedLines("test", deleteLine)

	fmt.Println(fc.Contents)

	require.Equal(t, 2, contentLen(fc))
}

func TestRemoveMultiLine(t *testing.T) {
	content := `SourceCode goes here!
// !spawntag:test
these
lines
are
removed
// !spawntag:test
final line`

	fc := &spawn.FileContent{
		Contents: content,
	}

	require.Equal(t, 8, contentLen(fc))

	deleteLine := true
	fc.RemoveTaggedLines("test", deleteLine)

	fmt.Println(fc.Contents)

	require.Equal(t, 2, contentLen(fc))
}

func TestLineEndsWithSymbol(t *testing.T) {
	require.True(t, spawn.DoesLineEndWithOpenSymbol(`tokenfactory.NewKeeper(`))
	require.True(t, spawn.DoesLineEndWithOpenSymbol(`tokenfactory.NewKeeper( // comment`))
	require.True(t, spawn.DoesLineEndWithOpenSymbol(`tokenfactory.NewKeeper(    `))
	require.True(t, spawn.DoesLineEndWithOpenSymbol(`tokenfactory{`))
	require.True(t, spawn.DoesLineEndWithOpenSymbol(`tokenfactory{ // comment}`))

	require.False(t, spawn.DoesLineEndWithOpenSymbol(` ) `))
	require.False(t, spawn.DoesLineEndWithOpenSymbol(` )((((())))(}{}`))
}

func contentLen(fs *spawn.FileContent) int {
	return len(strings.Split(fs.Contents, "\n"))
}

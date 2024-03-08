package spawn

import (
	"log/slog"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUncomment(t *testing.T) {
	content := `SourceCode goes here!
//this is line 1 which should be uncommented // ?spawntag:test
untouched line 3 // spawntag:test`

	fc := &FileContent{
		Contents: content,
	}

	fc.HandleCommentSwaps("test")

	// the line should be uncommented
	require.Equal(t, "this is line 1 which should be uncommented", strings.Split(fc.Contents, "\n")[1])
}

func TestRemoveLineWithTag(t *testing.T) {
	fc := &FileContent{
		Contents: `SourceCode goes here!
		//this is line 1 which should be uncommented // ?spawntag::test
		this line gets deleted //spawntag:test`,
	}

	deleteLine := true
	fc.RemoveTaggedLines("test", deleteLine)
	require.Equal(t, 2, contentLen(fc))
}

func TestRemoveLeftOverComments(t *testing.T) {
	fc := &FileContent{
		Contents: `SourceCode goes here!
		this line stays //spawntag:test
		this line stays too // ?spawntag:testing
		// <spawntag:multiline
		this line stays too
		// spawntag:multiline>
`,
	}

	deleteLine := false
	fc.RemoveTaggedLines("", deleteLine)

	require.False(t, strings.Contains(fc.Contents, "spawntag"), fc.Contents)
}

func TestBatchRemoveText(t *testing.T) {
	happyDir := path.Join("app", "app.go")

	fc := &FileContent{
		RelativePath: happyDir,
		Contents: `// Some random comment here
		first line of source
		second line of source

		app.TestKeeper = testkeeper.NewKeeper(
			appCodec,
			app.keys[testtypes.StoreKey],
			app.AccountKeeper,
			app.BankKeeper,
			app.DistrKeeper,
			[]string{
				testtypes.1,
				testtypes.2,
				testtypes.3,
			},
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		)
		third line of source
		// some comment`,
		Logger: slog.Default(),
	}

	fc.RemoveModuleFromText("test", path.Join("wrong-dir", "app.go"))
	require.Equal(t, 19, contentLen(fc))

	fc.RemoveModuleFromText("test", happyDir)
	require.Equal(t, 6, contentLen(fc))
}

func TestRemoveMultiLine(t *testing.T) {
	content := `SourceCode goes here!
// <spawntag:test
these
lines
are
removed
// spawntag:test>
final line`

	fc := &FileContent{
		Contents: content,
		Logger:   slog.Default(),
	}

	require.Equal(t, 8, contentLen(fc))

	deleteLine := true
	fc.RemoveTaggedLines("test", deleteLine)

	require.Equal(t, 2, contentLen(fc))
}

func TestCommentText(t *testing.T) {
	require.Equal(t, getCommentText("test // my comment"), "my comment")
	require.Equal(t, getCommentText("test //my comment"), "my comment")
}

func TestLineEndsWithSymbol(t *testing.T) {
	require.True(t, DoesLineEndWithOpenSymbol(`tokenfactory.NewKeeper(`))
	require.True(t, DoesLineEndWithOpenSymbol(`tokenfactory.NewKeeper( // comment`))
	require.True(t, DoesLineEndWithOpenSymbol(`tokenfactory.NewKeeper(    `))
	require.True(t, DoesLineEndWithOpenSymbol(`tokenfactory{`))
	require.True(t, DoesLineEndWithOpenSymbol(`tokenfactory{ // comment}`))

	require.False(t, DoesLineEndWithOpenSymbol(` ) `))
	require.False(t, DoesLineEndWithOpenSymbol(` )((((())))(}{}`))
}

func contentLen(fs *FileContent) int {
	return len(strings.Split(fs.Contents, "\n"))
}

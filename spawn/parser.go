package spawn

import (
	"fmt"
	"strings"
)

const (
	ExpectedFormat    = "// spawntag:"
	CommentSwapFormat = "?spawntag:"
)

func hasIgnoreComment(line string) bool {
	return strings.Contains(line, "//ignore") || strings.Contains(line, "// ignore") || strings.Contains(line, "spawntag:ignore")
}

// Sometimes we remove a module line and would like to swap it for another.
func (fc *FileContent) HandleCommentSwaps(name string) {
	if !strings.Contains(fc.Contents, CommentSwapFormat) {
		return
	}

	splitContent := strings.Split(fc.Contents, "\n")

	newContent := make([]string, 0, len(splitContent))

	uncomment := fmt.Sprintf("%s:%s", CommentSwapFormat, name)

	for idx, line := range splitContent {
		hasUncommentTag := strings.Contains(line, uncomment)

		if hasUncommentTag {
			line = strings.Replace(line, "//", "", 1)
			line = strings.TrimRight(strings.Replace(line, fmt.Sprintf("// %s", uncomment), "", 1), " ")
			fmt.Printf("uncomment %s: %d, %s\n", name, idx, line)
		}

		newContent = append(newContent, line)
	}

	fc.Contents = strings.Join(newContent, "\n")
}

// RemoveTaggedLines deletes tagged lines or just removes the comment if desired.
func (fc *FileContent) RemoveTaggedLines(name string, deleteLine bool) {
	splitContent := strings.Split(fc.Contents, "\n")
	newContent := make([]string, 0, len(splitContent))

	startIdx := -1
	for idx, line := range splitContent {
		hasTag := strings.Contains(line, fmt.Sprintf("spawntag:%s", name))
		hasMultiLineTag := strings.Contains(line, fmt.Sprintf("!spawntag:%s", name))

		// if the line has a tag, and the tag starts with a !, then we will continue until we
		// find the end of the tag with another.
		if startIdx != -1 {
			if !hasMultiLineTag {
				continue
			}

			startIdx = -1
			fmt.Println("endIdx:", idx, line)
			continue
		}

		if hasMultiLineTag {
			if !deleteLine {
				continue
			}

			startIdx = idx
			fmt.Printf("startIdx %s: %d, %s\n", name, idx, line)
			continue
		}

		if hasTag {
			if deleteLine {
				continue
			}

			line = removeJustSpawnTagLineComment(line)
		}

		newContent = append(newContent, line)
	}

	// return []byte(strings.Join(newContent, "\n"))
	fc.Contents = strings.Join(newContent, "\n")
}

// removeLineComment removes just the spawntag comment from a line of code.
// this way it is not user facing
func removeJustSpawnTagLineComment(line string) string {
	// QOL for us to not tear our hair out if we have a space or not
	// Could do this for all contents on load?
	line = strings.ReplaceAll(line, "//spawntag:", ExpectedFormat)

	line = strings.Split(line, ExpectedFormat)[0]
	return strings.TrimRight(line, " ")
}

// RemoveGeneralModule removes any matching names from the fileContent.
// i.e. if moduleFind is "tokenfactory" any lines with "tokenfactory" will be removed
// including comments.
// If an import or other line depends on a solo module a user wishes to remove, add a comment to the line
// such as `// spawntag:tokenfactory` to also remove other lines within the simapp template
func (fc *FileContent) RemoveModuleFromText(removeText string, pathSuffix ...string) {
	if !fc.InPaths(pathSuffix) {
		return
	}

	splitContent := strings.Split(fc.Contents, "\n")
	newContent := make([]string, 0, len(splitContent))

	startIdx := -1
	for idx, line := range splitContent {
		// if we are in a startIdx, then we need to continue until we find the close parenthesis (i.e. NewKeeper)
		if startIdx != -1 {
			fmt.Printf("rm %s startIdx: %d, %s\n", removeText, idx, line)
			if strings.TrimSpace(line) == ")" || strings.TrimSpace(line) == "}" {
				fmt.Println("endIdx:", idx, line)
				startIdx = -1
				continue
			}

			continue
		}

		// if line contains //ignore or // ignore, then we use that line
		// useful if some text is 'wasm' as a bech32 prefix, not a variable / type we need to remove.
		if hasIgnoreComment(line) {
			fmt.Printf("Ignoring removal: %s: %d, %s\n", removeText, idx, line)
			newContent = append(newContent, line)
			continue
		}

		lineHas := strings.Contains(line, removeText)

		if lineHas && DoesLineEndWithOpenSymbol(line) {
			startIdx = idx
			fmt.Printf("startIdx %s: %d, %s\n", removeText, idx, line)
			continue
		}

		if lineHas {
			fmt.Printf("rm %s: %d, %s\n", removeText, idx, line)
			continue
		}

		newContent = append(newContent, line)
	}

	fc.Contents = strings.Join(newContent, "\n")
}

// doesLineEndWithOpenSymbol returns true if the end of a line opens a statement such as a multi-line function.
func DoesLineEndWithOpenSymbol(line string) bool {
	// remove comment if there is one
	if strings.Contains(line, "//") {
		line = strings.Split(line, "//")[0]
	}

	return strings.HasSuffix(strings.TrimSpace(line), "(") || strings.HasSuffix(strings.TrimSpace(line), "{")
}

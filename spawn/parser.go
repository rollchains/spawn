package spawn

import (
	"fmt"
	"strings"
)

const (
	// StdFormat is the standard format for removing a line if a feature is removed.
	StdFormat = "spawntag:%s"

	// ExpectedFormat is the standard format for removing a line if a module is removed.
	// e.g. // spawntag:tokenfactory would remove the line if tokenfactory is removed.
	// NOTE: This is not user facing, and is only used for internal parsing of the simapp.
	ExpectedFormat = "// spawntag:"

	// CommentSwapFormat is the format for swapping a line with another if a module is removed.
	CommentSwapFormat = "?spawntag:%s"

	// MultiLineStartFormat is the format for starting a multi-line comment which removes all text
	// until the end of the comment.
	// <spawntag:[searchTerm]
	MultiLineStartFormat = "<" + StdFormat

	// spawntag:[searchTerm]>
	MultiLineEndFormat = StdFormat + ">"
)

// Sometimes we remove a module line and would like to swap it for another.
func (fc *FileContent) HandleCommentSwaps(name string) {
	splitContent := strings.Split(fc.Contents, "\n")

	uncommentSearchTag := fmt.Sprintf(CommentSwapFormat, name)

	for idx, line := range splitContent {
		// If the line does not have the comment swap tag, then continue
		if !strings.Contains(line, uncommentSearchTag) {
			continue
		}

		// removes the // spawntag:[name] comment from the end of the source code
		line = removeSpawnTagLineComment(line, uncommentSearchTag)

		// uncomments the line (to expose the source code)
		line = uncommentLine(line)

		// Since we are just uncommenting the line, it's safe to just replace the line at the index
		splitContent[idx] = line

	}

	fc.Contents = strings.Join(splitContent, "\n")
}

// RemoveTaggedLines deletes tagged lines or just removes the comment if desired.
func (fc *FileContent) RemoveTaggedLines(name string, deleteLine bool) {
	splitContent := strings.Split(fc.Contents, "\n")
	newContent := make([]string, 0, len(splitContent))

	startMultiLineDelete := false
	for idx, line := range splitContent {

		// if the line has a tag, and the tag starts with a !, then we will continue until we
		// find the end of the tag with another.
		if startMultiLineDelete {
			hasMultiLineEndTag := strings.Contains(line, fmt.Sprintf(MultiLineEndFormat, name))
			if !hasMultiLineEndTag {
				continue
			}

			// the line which has the closing multiline end tag, we then continue to add lines as normal
			startMultiLineDelete = false
			fmt.Println("endIdx:", idx, line)
			continue
		}

		// <spawntag:[searchTerm]
		if strings.Contains(line, fmt.Sprintf(MultiLineStartFormat, name)) {
			if !deleteLine {
				continue
			}

			startMultiLineDelete = true
			fmt.Printf("startIdx %s: %d, %s\n", name, idx, line)
			continue
		}

		// remove a line if it contains spawntag:[searchTerm]
		if strings.Contains(line, fmt.Sprintf(StdFormat, name)) {
			if deleteLine {
				continue
			}

			line = removeSpawnTagLineComment(line, ExpectedFormat)
		}

		newContent = append(newContent, line)
	}

	// return []byte(strings.Join(newContent, "\n"))
	fc.Contents = strings.Join(newContent, "\n")
}

// removeSpawnTagLineComment removes just the spawntag comment from a line of code.
func removeSpawnTagLineComment(line string, tag string) string {
	// QOL for us to not tear our hair out if we have a space or not
	// Could do this for all contents on load?
	line = strings.ReplaceAll(line, "//spawntag:", ExpectedFormat)

	line = strings.Split(line, fmt.Sprintf("// %s", tag))[0]
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

		// if line contains //spawntag:ignorem then we use that line
		// useful if some text is 'wasm' as a bech32 prefix, not a variable / type we need to remove.
		if strings.Contains(line, fmt.Sprintf(StdFormat, "ignore")) {
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

// getCommentText returns the trimmed text from a line comment.
func getCommentText(line string) string {
	if strings.Contains(line, "//") {
		text := strings.Split(line, "//")[1]
		return strings.TrimSpace(text)
	}

	return ""
}

func uncommentLine(line string) string {
	return strings.Replace(line, "//", "", 1)
}

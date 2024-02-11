package spawn

import (
	"fmt"
	"strings"
)

const expectedFormat = "// spawntag:"

// Sometimes if we remove a module, we want to delete one line and use another.
func (fc *FileContent) HandleCommentSwaps(name string) {
	newContent := make([]string, 0, len(strings.Split(fc.Contents, "\n")))

	uncomment := fmt.Sprintf("?spawntag:%s", name)

	for idx, line := range strings.Split(fc.Contents, "\n") {
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
	newContent := make([]string, 0, len(strings.Split(fc.Contents, "\n")))

	startIdx := -1
	for idx, line := range strings.Split(fc.Contents, "\n") {
		// TODO: regex anything in between // and spawntag such as spaces, symbols, etc?
		// TODO: Do this for all content on load?
		line = strings.ReplaceAll(line, "//spawntag:", expectedFormat) // just QOL for us to not tear our hair out

		hasTag := strings.Contains(line, fmt.Sprintf("spawntag:%s", name))
		hasMultiLineTag := strings.Contains(line, fmt.Sprintf("!spawntag:%s", name))

		// if the line has a tag, and the tag starts with a !, then we will continue until we find the end of the tag with another.
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

			line = strings.Split(line, expectedFormat)[0]
			line = strings.TrimRight(line, " ")
		}

		newContent = append(newContent, line)
	}

	// return []byte(strings.Join(newContent, "\n"))
	fc.Contents = strings.Join(newContent, "\n")
}

func (fc *FileContent) DeleteContents(path string) {
	if strings.HasSuffix(fc.NewPath, path) {
		fmt.Println("Deleting contents for", path)
		fc.Contents = ""
	}
}

// RemoveGeneralModule removes any matching names from the fileContent.
// i.e. if moduleFind is "tokenfactory" any lines with "tokenfactory" will be removed
// including comments.
// If an import or other line depends on a solo module a user wishes to remove, add a comment to the line
// such as `// tag:tokenfactory` to also remove other lines within the simapp template
func (fc *FileContent) RemoveModuleFromText(removeText string, pathSuffix ...string) {
	// if !strings.HasSuffix(fc.NewPath, pathSuffix) {
	// 	return
	// }

	found := false
	for _, suffix := range pathSuffix {
		if strings.HasSuffix(fc.RelativePath, suffix) {
			found = true
			break
		}
	}
	if !found {
		return
	}

	newContent := make([]string, 0, len(strings.Split(fc.Contents, "\n")))

	startIdx := -1
	for idx, line := range strings.Split(fc.Contents, "\n") {
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

		lineHas := strings.Contains(line, removeText)

		// if line contains //ignore or // ignore, then we use that line
		// useful if some text is 'wasm' as a bech32 prefix, not a variable / type.
		if strings.Contains(line, "//ignore") || strings.Contains(line, "// ignore") {
			fmt.Printf("Ignoring removal: %s: %d, %s\n", removeText, idx, line)
			newContent = append(newContent, line)
			continue
		}

		if lineHas && (strings.HasSuffix(strings.TrimSpace(line), "(") || strings.HasSuffix(strings.TrimSpace(line), "{")) {
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

// given a go mod, remove a line within the file content
func RemoveGoModImport(module string, fileContent []byte) []byte {
	fcs := string(fileContent)
	lines := strings.Split(fcs, "\n")

	newLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.Contains(line, module) {
			newLines = append(newLines, line)
		}
	}

	return []byte(strings.Join(newLines, "\n"))
}

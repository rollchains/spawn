package spawn

import (
	"fmt"
	"os"
	"path"
	"strings"

	modfile "golang.org/x/mod/modfile"
)

// FindLinesWithText returns the start and end index of a block of text in a slice of strings.
// This starts with the text searched for, and continues to look until and end bracket or parenthesis is found.
// This allows for the finding of multi-line function signatures, if-else blocks, etc.
func FindLinesWithText(src []string, text string) (startIdx, endIdx int) {
	startMultiLineFind := false
	for idx, line := range src {
		if startMultiLineFind {
			if strings.TrimSpace(line) == ")" || strings.TrimSpace(line) == "}" {
				return startIdx, idx + 1
			}
		}

		if strings.Contains(line, text) {
			startMultiLineFind = true
			startIdx = idx
			continue
		}
	}

	return 0, 0
}

// FindLineWithText returns the index of a line in a slice of strings that contains the given text.
func FindLineWithText(src []string, text string) (lineNum int) {
	for i, line := range src {
		if strings.Contains(line, text) {
			return i
		}
	}

	return 0
}

// ParseFileImports reads the content of a file and returns the import strings, the start and end line numbers of the import block.
// It starts reading at `import (` and stops at the first `)`. in the file.
func ParseFileImports(filePath string) ([]string, int, int, error) {
	// Read the content of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, 0, 0, err
	}

	// Split the content into lines
	lines := strings.Split(string(content), "\n")

	// Find the import block and its boundaries
	importStartLine := -1
	importEndLine := -1
	for i, line := range lines {
		if strings.Contains(line, "import (") {
			importStartLine = i + 1 // Line numbers start from 1
		} else if importStartLine != -1 && strings.Contains(line, ")") {
			importEndLine = i + 1 // Line numbers start from 1
			break
		}
	}

	// If no import block found, return empty slice and line numbers as 0
	if importStartLine == -1 || importEndLine == -1 {
		return []string{}, 0, 0, nil
	}

	// Extract import strings within the import block
	var imports []string
	for _, line := range lines[importStartLine:importEndLine] {
		if strings.Contains(line, "\"") {
			imports = append(imports, line)
		}
	}

	return imports, importStartLine, importEndLine, nil
}

// ReadCurrentModuleName reads the go package module name from the go.mod file on the host machine.
func ReadCurrentGoModuleName(loc string) string {
	if !strings.HasSuffix(loc, "go.mod") {
		loc = path.Join(loc, "go.mod")
	}

	goModBz, err := os.ReadFile(loc)
	if err != nil {
		fmt.Println("Error reading go.mod file: ", err)
		return ""
	}

	return modfile.ModulePath(goModBz)
}

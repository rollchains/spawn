package spawn

import (
	"fmt"
	"os"
	"path"

	modfile "golang.org/x/mod/modfile"
)

func getModPath() string {
	// used when you `spawn new`
	goModPath := path.Join("simapp", "go.mod")

	// testing mode:
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		// specific unit test run
		goModPath = path.Join("..", "simapp", "go.mod")
		// go test ./...
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			goModPath = path.Join("..", "..", "simapp", "go.mod")
		}
	}

	return goModPath
}

// ParseVersionFromGoMod parses out the versions for a given goPath
// Ex: ParseVersionFromGoMod("github.com/cosmos/cosmos-sdk", false) returns v0.50.X
// Ex: ParseVersionFromGoMod("github.com/cosmos/cosmos-sdk", true) returns 0.50.X
func ParseVersionFromGoMod(goPath string, removePrefixedV bool) (string, error) {
	goModPath := getModPath()

	c, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("error reading go.mod file: %w", err)
	}

	f, err := modfile.Parse(goModPath, c, nil)
	if err != nil {
		return "", fmt.Errorf("error parsing go.mod file: %w", err)
	}

	for _, r := range f.Require {
		if r.Mod.Path == goPath {
			v := r.Mod.Version
			if removePrefixedV && len(v) > 0 && v[0] == 'v' {
				v = v[1:]
			}
			return v, nil
		}
	}

	// no error if not found, we just return nothing
	return "", fmt.Errorf("module %s not found in go.mod", goPath)
}

func MustParseVersionFromGoMod(goPath string, removePrefixedV bool) string {
	v, err := ParseVersionFromGoMod(goPath, removePrefixedV)
	if err != nil {
		panic(err)
	}
	return v
}

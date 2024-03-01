package spawn

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

// The Name, RequestMsg, and ResponseMsg fields of rpc services
type ProtoRPC struct {
	// The name of the proto RPC service (i.e. rpc Params would be Params for the name)
	Name string
	// The request object, such as QueryParamsRequest (queries) or MsgUpdateParams (txs)
	Req string
	// The response object, such as QueryParamsResponse (queries) or MsgUpdateParamsResponse (txs)
	Res string

	// the relative directory location this proto file is location (x/mymodule/types)
	Location string
	// The name of the module
	Module string
	// The type of file this proto service is
	FType FileType
}

// ProtoServiceParser parses out a proto file and returns all the services within it.
func ProtoServiceParser(content []byte, pkgDir string, ft FileType) []*ProtoRPC {
	qss := make([]*ProtoRPC, 0)
	c := strings.Split(string(content), "\n")

	for idx, line := range c {
		if strings.Contains(line, "rpc ") {
			fmt.Println("Found rpc line: ", strings.Trim(line, " "))

			// if line does not end with {, we also need to load the next line
			if !strings.HasSuffix(line, "{") {
				line = line + c[idx+1]
			}

			line = strings.Trim(line, " ")

			line = strings.NewReplacer("rpc", "", "returns", "", "(", " ", ")", " ", "{", "", "}", "").Replace(line)

			words := strings.Fields(line)
			qss = append(qss, &ProtoRPC{
				Name:     words[0],
				Req:      words[1],
				Res:      words[2],
				Location: pkgDir,
				FType:    ft,
			})
		}
	}

	return qss
}

// FileType tells the application which type of proto file is it so we can sort Txs from Queries
type FileType string

const (
	Tx    FileType = "tx"
	Query FileType = "query"
	None  FileType = "none"
)

// returns "tx" or "query" depending on the content of the file
func SortContentToFileType(bz []byte) FileType {
	res := string(bz)

	// if `service Query` or `message Query` found in the file, it's a query
	if strings.Contains(res, "service Query") || strings.Contains(res, "message Query") {
		return Query
	}

	// if `service Msg` or `service Tx` or `message Msg`
	if strings.Contains(res, "service Msg") || strings.Contains(res, "service Tx") || strings.Contains(res, "message Msg") {
		return Tx
	}

	return None
}

// GetGoPackageLocationOfFiles parses the proto content pulling out the relative path
// of the go package location.
// option go_package = "github.com/rollchains/mychain/x/cnd/types"; -> x/cnd/types
func GetGoPackageLocationOfFiles(bz []byte) string {
	modName := ReadCurrentGoModuleName("go.mod")

	for _, line := range strings.Split(string(bz), "\n") {
		if strings.Contains(line, "option go_package") {
			// option go_package = "github.com/rollchains/mychain/x/cnd/types";
			line = strings.Trim(line, " ")

			// line = strings.NewReplacer("option go_package", "", "=", "", ";", "", , "", "\"", "").Replace(line)

			// x/cnd/types";
			line = strings.Split(line, fmt.Sprintf("%s/", modName))[1]
			// x/cnd/types
			line = strings.Split(line, "\";")[0]

			return strings.Trim(line, " ")
		}
	}

	return ""
}

// helpers

/*
 TODO: is this used or needed? (was at the top of rthe proto service generator)
func GetProtoDirectories(protoAbsPath string, args ...string) []string {
	dirs, err := os.ReadDir(protoAbsPath)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	absDirs := make([]string, 0)
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		if len(args) > 0 && dir.Name() != args[0] {
			continue
		}

		absDirs = append(absDirs, path.Join(protoAbsPath, dir.Name()))
	}

	fmt.Println("Found dirs: ", absDirs)

	return absDirs
}
*/

// Converts .proto files into a mapping depending on the type.
// TODO: is the 2nd map of FileType required since ProtoRPC has it anyways?
func GetModuleMapFromProto(absProtoPath string) map[string][]*ProtoRPC {
	modules := make(map[string][]*ProtoRPC)

	fs.WalkDir(os.DirFS(absProtoPath), ".", func(relPath string, d fs.DirEntry, e error) error {
		if !strings.HasSuffix(relPath, ".proto") {
			return nil
		}

		// read file content
		content, err := os.ReadFile(path.Join(absProtoPath, relPath))
		if err != nil {
			fmt.Println("Error: ", err)
		}

		fileType := SortContentToFileType(content)

		parent := path.Dir(relPath)
		parent = strings.Split(parent, "/")[0]

		// add/append to modules
		if _, ok := modules[parent]; !ok {
			modules[parent] = make([]*ProtoRPC, 0)
		}

		goPkgDir := GetGoPackageLocationOfFiles(content)

		switch fileType {
		case Tx:
			fmt.Println("File is a transaction")
			tx := ProtoServiceParser(content, goPkgDir, Tx)
			modules[parent] = append(modules[parent], tx...)

		case Query:
			fmt.Println("File is a query")
			query := ProtoServiceParser(content, goPkgDir, Query)
			// modules[parent][Query] = append(modules[parent][Query], query...)
			modules[parent] = append(modules[parent], query...)
		case None:
			fmt.Println("File is neither a transaction nor a query")
		}

		return nil
	})

	fmt.Printf("Modules: %+v\n", modules)
	return modules
}

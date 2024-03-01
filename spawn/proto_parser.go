package spawn

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

// A Proto server RPC method.
type ProtoRPC struct {
	// The name of the proto RPC service (i.e. rpc Params would be Params for the name)
	Name string
	// The request object, such as QueryParamsRequest (queries) or MsgUpdateParams (txs)
	Req string
	// The response object, such as QueryParamsResponse (queries) or MsgUpdateParamsResponse (txs)
	Res string

	// The name of the module
	Module string
	// the relative directory location this proto file is location (x/mymodule/types)
	Location string
	// The type of file this proto service is
	FType FileType
	// Where there types.(Query/Msg)Server is located
	FileLoc string
}

// ProtoServiceParser parses out a proto file content and returns all the services within it.
func ProtoServiceParser(content []byte, pkgDir string, ft FileType) []*ProtoRPC {
	pRPCs := make([]*ProtoRPC, 0)
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
			pRPCs = append(pRPCs, &ProtoRPC{
				Name:     words[0],
				Req:      words[1],
				Res:      words[2],
				Location: pkgDir,
				FType:    ft,
			})
		}
	}

	return pRPCs
}

// FileType tells the application which type of proto file is it so we can sort Txs from Queries
type FileType string

const (
	Tx    FileType = "tx"
	Query FileType = "query"
	None  FileType = "none"
)

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

		fileType := GetFileTypeFromProtoContent(content)

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

// returns "tx" or "query" depending on the content of the file
func GetFileTypeFromProtoContent(bz []byte) FileType {
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

// type Modules struct {
// 	Modules map[string][]*ProtoRPC
// }

type MissingModules struct {
	// module name -> RPC methods
	Missing map[string][]*ProtoRPC
}

func GetCurrentRPCMethodsFromModuleProto(cwd string) *MissingModules {
	protoPath := path.Join(cwd, "proto")

	modules := GetModuleMapFromProto(protoPath)

	missing := make(map[string][]*ProtoRPC, 0)

	// TODO: currently if using multiple modules, it will run the code 2 times for generating missing methods
	for name, rpcMethods := range modules {
		fmt.Println("\n------------- Module: ", name)

		modulePath := path.Join(cwd, "x", name, "keeper") // hardcode for keeper is less than ideal, but will do for now

		// currentMethods := make(map[FileType][]string, 0) // tx/query -> methods
		txMethods := make([]string, 0)
		queryMethods := make([]string, 0)

		// msgServerFile := ""
		// queryServerFile := ""

		for _, rpc := range rpcMethods {
			rpc := rpc
			rpc.Module = name

			fmt.Println("\nService: ", rpc)

			// get files in service.Location
			files, err := os.ReadDir(modulePath)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			for _, f := range files {
				if strings.HasSuffix(f.Name(), "_test.go") {
					continue
				}

				content, err := os.ReadFile(path.Join(modulePath, f.Name()))
				if err != nil {
					fmt.Println("Error: ", err)
				}

				// if the file type is not the expected, continue
				// if the content of this file is not the same as the service we are tying to use, continue
				if rpc.FType != isFileQueryOrMsgServer(content) {
					continue
				}

				fmt.Println(" = File: ", f.Name())

				// Set the file type tot his file
				rpc.FileLoc = path.Join(modulePath, f.Name())

				switch rpc.FType {
				case Tx:
					// msgServerFile = path.Join(modulePath, f.Name())
				case Query:
					// queryServerFile = path.Join(modulePath, f.Name())
					rpc.FileLoc = path.Join(modulePath, f.Name())
				default:
					fmt.Println("Error: ", "Unknown FileType")
					panic("RUT ROE RAGGY")
				}

				// find any line with `func ` in it

				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					// receiver func
					if strings.Contains(line, "func (") {
						if strings.Contains(line, rpc.Name) {
							switch rpc.FType {
							case Tx:
								txMethods = append(txMethods, parseReceiverMethodName(line))
							case Query:
								queryMethods = append(queryMethods, parseReceiverMethodName(line))
							default:
								fmt.Println("Error: ", "Unknown FileType")
								panic("RUT ROE RAGGY")
							}
						}
					}
				}
			}

			// filePaths[name] = map[FileType]string{
			// 	Tx:    msgServerFile,
			// 	Query: queryServerFile,
			// }
		}

		// map[query:[Params] tx:[UpdateParams]]
		fmt.Println("\n-------- Current Tx Methods: ", txMethods)
		fmt.Println("-------- Current Query Methods: ", queryMethods)

		// iterate services again and apply any missing methods to the file
		for name, rpcs := range modules {
			for _, rpc := range rpcs {
				// current := currentMethods[fileType]

				var current []string
				switch rpc.FType {
				case Tx:
					current = txMethods
				case Query:
					current = queryMethods
				default:
					fmt.Println("Error: ", "Unknown FileType")
					panic("RUT ROE RAGGY")
				}

				fmt.Println("   Current: ", rpc.FType, current)

				found := false
				for _, method := range current {
					method := method
					if method == rpc.Name {
						found = true
						break
					}
				}

				if found {
					fmt.Println("  - Found: ", rpc.Name)
					continue
				}

				if _, ok := missing[name]; !ok {
					missing[name] = make([]*ProtoRPC, 0)
				}

				// MISSING METHOD
				// fmt.Println("Missing method: ", service.Name, service.Req, service.Res, service.Location)
				missing[name] = append(missing[name], rpc)
				fmt.Println("  - Missing: ", rpc.Name)
			}
		}

		fmt.Println("\n\nMissing: ")
		for name, rpcs := range missing {
			fmt.Println("  - Module: ", name)
			for _, rpc := range rpcs {
				fmt.Println("    - ", rpc.Name, rpc.Req, rpc.Res)
			}
		}
	}

	return &MissingModules{
		Missing: missing,
	}
}

func parseReceiverMethodName(f string) string {
	// given a string of text like `func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {`
	// parse out Params, req and the response

	name := ""

	f = strings.ReplaceAll(f, "func (", "")
	parts := strings.Split(f, ") ")
	name = strings.Split(parts[1], "(")[0]

	return strings.Trim(name, " ")
}

func isFileQueryOrMsgServer(bz []byte) FileType {
	s := strings.ToLower(string(bz))

	if strings.Contains(s, "queryserver") || strings.Contains(s, "querier") {
		return Query
	}

	if strings.Contains(s, "msgserver") || strings.Contains(s, "msgservice") {
		return Tx
	}

	return None
}

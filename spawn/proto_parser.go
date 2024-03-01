package spawn

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

// BuildProtoInterfaceStub builds the stub for the proto interface depending on the ProtoRPC's file type.
func (pr ProtoRPC) BuildProtoInterfaceStub() string {
	if pr.FType == Tx {
		return fmt.Sprintf(`// %s implements types.MsgServer.
func (ms msgServer) %s(ctx context.Context, msg *types.%s) (*types.%s, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	panic("%s is unimplemented")
	return &types.%s{}, nil
}
`, pr.Name, pr.Name, pr.Req, pr.Res, pr.Name, pr.Res)
	} else if pr.FType == Query {
		return fmt.Sprintf(`// %s implements types.QueryServer.
func (k Querier) %s(goCtx context.Context, req *types.%s) (*types.%s, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	panic("%s is unimplemented")
	return &types.%s{}, nil
}
`, pr.Name, pr.Name, pr.Req, pr.Res, pr.Name, pr.Res)
	} else {
		panic("Unknown FileType for: " + pr.Name)
	}
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

// Converts .proto files into a mapping depending on the type.
func GetCurrentModuleRPCsFromProto(absProtoPath string) ModuleMapping {
	modules := make(ModuleMapping)

	fs.WalkDir(os.DirFS(absProtoPath), ".", func(relPath string, d fs.DirEntry, e error) error {
		if !strings.HasSuffix(relPath, ".proto") {
			return nil
		}

		content, err := os.ReadFile(path.Join(absProtoPath, relPath))
		if err != nil {
			fmt.Println("Error: ", err)
		}

		fileType := GetFileTypeFromProtoContent(content)
		if fileType == None {
			return nil
		}

		goPkgDir := GetGoPackageLocationOfFiles(content)

		rpcs := ProtoServiceParser(content, goPkgDir, fileType)

		parent := path.Dir(relPath)
		parent = strings.Split(parent, "/")[0]

		if _, ok := modules[parent]; !ok {
			modules[parent] = make([]*ProtoRPC, 0)
		}

		modules[parent] = append(modules[parent], rpcs...)

		return nil
	})

	modules.Print()

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

func GetMissingRPCMethodsFromModuleProto(cwd string) (ModuleMapping, error) {
	protoPath := path.Join(cwd, "proto")
	modules := GetCurrentModuleRPCsFromProto(protoPath)

	missing := make(ModuleMapping, 0)

	// txMethods := make([]string, 0)
	// 	queryMethods := make([]string, 0)

	txMethods := make(map[string][]string)
	queryMethods := make(map[string][]string)

	for name, rpcMethods := range modules {
		fmt.Println("\n------------- Module: ", name)

		modulePath := path.Join(cwd, "x", name, "keeper") // hardcode for keeper is less than ideal, but will do for now

		for _, rpc := range rpcMethods {
			rpc := rpc
			rpc.Module = name

			fmt.Println("\nService: ", rpc)

			// get files in service.Location
			files, err := os.ReadDir(modulePath)
			if err != nil {
				return nil, err
			}

			for _, f := range files {
				if strings.HasSuffix(f.Name(), "_test.go") {
					continue
				}

				content, err := os.ReadFile(path.Join(modulePath, f.Name()))
				if err != nil {
					// fmt.Println("Error: ", err)
					return nil, err
				}

				// if the file type is not the expected, continue
				// if the content of this file is not the same as the service we are tying to use, continue
				if rpc.FType != isFileQueryOrMsgServer(content) {
					continue
				}

				fmt.Println(" = File: ", f.Name())

				// Set the file type tot his file
				rpc.FileLoc = path.Join(modulePath, f.Name())

				// find any line with `func ` in it
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					// receiver func
					if strings.Contains(line, "func (") {
						if strings.Contains(line, rpc.Name+"(") {
							switch rpc.FType {
							case Tx:
								txMethods[name] = append(txMethods[name], parseReceiverMethodName(line))
							case Query:
								queryMethods[name] = append(queryMethods[name], parseReceiverMethodName(line))
							default:
								fmt.Println("Error: ", "Unknown FileType")
								panic("RUT ROE RAGGY")
							}
						}
					}
				}
			}
		}

		// map[query:[Params] tx:[UpdateParams]]
		fmt.Println("\n-------- Current Tx Methods: ", txMethods)
		fmt.Println("-------- Current Query Methods: ", queryMethods)
	}

	// iterate services again and apply any missing methods to the file
	for name, rpcs := range modules {
		rpcs := rpcs

		if _, ok := missing[name]; !ok {
			missing[name] = make([]*ProtoRPC, 0)
		}

		for _, rpc := range rpcs {
			rpc := rpc
			// current := currentMethods[fileType]

			var current []string
			switch rpc.FType {
			case Tx:
				current = txMethods[name]
			case Query:
				current = queryMethods[name]
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

			alreadyIncluded := false
			for _, m := range missing[name] {
				if m.Name == rpc.Name {
					alreadyIncluded = true
					break
				}
			}
			if alreadyIncluded {
				continue
			}

			// MISSING METHOD
			// fmt.Println("Missing method: ", service.Name, service.Req, service.Res, service.Location)
			missing[name] = append(missing[name], rpc)
			fmt.Println("  - Missing: ", rpc.Name, name)
		}
	}

	fmt.Println("\n\nMissing: ")
	for name, rpcs := range missing {
		fmt.Println("  - Module: ", name)
		for _, rpc := range rpcs {
			fmt.Println("    - ", rpc.Name, rpc.Req, rpc.Res)
		}
	}

	return missing, nil
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

func ApplyMissingRPCMethodsToGoSourceFiles(missingRPCMethods ModuleMapping) error {
	for _, missing := range missingRPCMethods {
		for _, miss := range missing {
			miss := miss
			fmt.Println("Module: ", miss.Module, "FType: ", miss.FType)

			fileLoc := miss.FileLoc
			fmt.Println("File: ", fileLoc)

			content, err := os.ReadFile(fileLoc)
			if err != nil {
				return fmt.Errorf("error: %s, file: %s", err.Error(), fileLoc)
			}

			fmt.Println("Append to file: ", miss.FType, miss.Name, miss.Req, miss.Res)

			code := miss.BuildProtoInterfaceStub()
			if len(code) == 0 {
				continue
			}

			// append to the file content after a new line at the end
			content = append(content, []byte("\n"+code)...)

			if err := os.WriteFile(fileLoc, content, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

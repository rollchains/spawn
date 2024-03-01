package spawn

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"strings"
)

// BuildProtoInterfaceStub returns the string to save to the file for the msgServer or Querier.
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

// ProtoServiceParser parses out a proto file content and returns all the RPC services within it.
func ProtoServiceParser(logger *slog.Logger, content []byte, pkgDir string, ft FileType, fileLoc string) []*ProtoRPC {
	pRPCs := make([]*ProtoRPC, 0)
	c := strings.Split(string(content), "\n")

	moduleName := GetProtoPackageName(content)

	for idx, line := range c {
		line = strings.TrimLeft(line, " ")
		line = strings.TrimLeft(line, "\t")

		if strings.HasPrefix(line, "rpc ") {
			line = strings.Trim(line, " ")
			logger.Debug("proto file", "rpc line", line)

			// if line does not end with {, we also need to load the next line (multi line proto from linting)
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
				Module:   moduleName,
				FileLoc:  fileLoc,
			})
		}
	}

	return pRPCs
}

// GetProtoPackageName inputs proto file content, then parse out the package (cosmos module) name
// package cnd.v1; returns cnd as the name.
func GetProtoPackageName(content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.Trim(line, " ")
		line = strings.Trim(line, "\t")

		if strings.HasPrefix(line, "package ") {
			// package cnd.v1;
			line = strings.Trim(line, " ")

			// cnd.v1;
			line = strings.Split(line, "package ")[1]

			// split at the first .
			line = strings.Split(line, ".")[0]

			return strings.Trim(line, " ")
		}
	}

	return ""
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

			// x/cnd/types";
			line := strings.Split(line, fmt.Sprintf("%s/", modName))[1]

			// x/cnd/types
			line = strings.Split(line, "\";")[0]

			return strings.Trim(line, " ")
		}
	}

	return ""
}

// Converts .proto files into a mapping depending on the type.
func GetCurrentModuleRPCsFromProto(logger *slog.Logger, absProtoPath string) ModuleMapping {
	modules := make(ModuleMapping)

	fs.WalkDir(os.DirFS(absProtoPath), ".", func(relPath string, d fs.DirEntry, e error) error {
		if !strings.HasSuffix(relPath, ".proto") {
			return nil
		}

		loc := path.Join(absProtoPath, relPath)

		content, err := os.ReadFile(loc)
		if err != nil {
			logger.Error("Error", "error", err)
		}

		fileType := FileTypeFromProtoContent(content)
		if fileType == None {
			return nil
		}

		goPkgDir := GetGoPackageLocationOfFiles(content)

		rpcs := ProtoServiceParser(logger, content, goPkgDir, fileType, loc)

		parent := path.Dir(relPath)
		parent = strings.Split(parent, "/")[0]

		if _, ok := modules[parent]; !ok {
			modules[parent] = make([]*ProtoRPC, 0)
		}

		modules[parent] = append(modules[parent], rpcs...)

		return nil
	})

	modules.Print(logger)

	return modules
}

// returns "tx" or "query" depending on the content of the file
func FileTypeFromProtoContent(bz []byte) FileType {
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

func GetMissingRPCMethodsFromModuleProto(logger *slog.Logger, cwd string) (ModuleMapping, error) {
	protoPath := path.Join(cwd, "proto")
	modules := GetCurrentModuleRPCsFromProto(logger, protoPath)

	missing := make(ModuleMapping, 0)

	txMethods := make(map[string][]string)
	queryMethods := make(map[string][]string)

	for name, rpcMethods := range modules {
		logger.Debug("module", "module", name)

		modulePath := path.Join(cwd, "x", name, "keeper") // hardcode for keeper is less than ideal, but will do for now

		for _, rpc := range rpcMethods {
			rpc := rpc
			rpc.Module = name

			logger.Debug("rpc", "rpc", rpc)

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
					logger.Error("error", "err", err)
					return nil, err
				}

				// if the file type is not the expected, continue
				// if the content of this file is not the same as the service we are tying to use, continue
				if rpc.FType != getFileType(content) {
					continue
				}

				logger.Debug(" = file", "file", f.Name())

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
								logger.Error("error", "err", "Unknown FileType")
								return nil, fmt.Errorf("unknown file type")
							}
						}
					}
				}

			}
		}

		logger.Debug("Current Tx Methods: ", "txMethods", txMethods)
		logger.Debug("Current Query Methods: ", "queryMethods", queryMethods)
	}

	// iterate services again and apply any missing methods to the file
	for name, rpcs := range modules {
		rpcs := rpcs

		if _, ok := missing[name]; !ok {
			missing[name] = make([]*ProtoRPC, 0)
		}

		for _, rpc := range rpcs {
			rpc := rpc

			var current []string
			switch rpc.FType {
			case Tx:
				current = txMethods[name]
			case Query:
				current = queryMethods[name]
			default:
				logger.Error("error", "err", "Unknown FileType")
				return nil, fmt.Errorf("unknown file type")
			}

			logger.Debug("    Current: ", "current", current)

			found := false
			for _, method := range current {
				method := method
				if method == rpc.Name {
					found = true
					break
				}
			}

			if found {
				continue
			}
			logger.Debug("  - Not Found: ", "rpc", rpc.Name)

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
			missing[name] = append(missing[name], rpc)
			logger.Debug("  - Missing: ", "rpc", rpc.Name, "module", name)
		}
	}

	missing.Print(logger)

	return missing, nil
}

// given a string of text like `func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {`
// parse out Params, req and the response
func parseReceiverMethodName(f string) string {
	name := ""

	// f = func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {

	// k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	f = strings.ReplaceAll(f, "func (", "")

	// [`k Querier` , `Params(c context.Context, req *types.QueryParamsRequest`, `(*types.QueryParamsResponse, error) {` ]
	parts := strings.Split(f, ") ")

	// [`Params`, `c context.Context, req *types.QueryParamsRequest`]
	name = strings.Split(parts[1], "(")[0]

	// `Params`
	return strings.Trim(name, " ")
}

func getFileType(bz []byte) FileType {
	s := strings.ToLower(string(bz))

	if strings.Contains(s, "queryserver") || strings.Contains(s, "querier") {
		return Query
	}

	if strings.Contains(s, "msgserver") || strings.Contains(s, "msgservice") {
		return Tx
	}

	return None
}

// ApplyMissingRPCMethodsToGoSourceFiles builds the proto interface stubs and appends them to the file for missing methods.
// If .proto file contained an rpc method for `Params` and `Other` but only `Params` is found in the querier, then `Other` is generated, appended, and saved.
func ApplyMissingRPCMethodsToGoSourceFiles(logger *slog.Logger, missingRPCMethods ModuleMapping) error {
	for _, missing := range missingRPCMethods {
		for _, rpc := range missing {
			miss := rpc
			fileLoc := miss.FileLoc

			logger.Debug("rpc info", "module", miss.Module, "fileLoc", fileLoc, "name", miss.Name, "ftype", miss.FType.String())

			content, err := os.ReadFile(fileLoc)
			if err != nil {
				return fmt.Errorf("error: %s, file: %s", err.Error(), fileLoc)
			}

			logger.Debug("Append to file: ", miss.Name, miss.Req, miss.Res, miss.Location)

			code := miss.BuildProtoInterfaceStub()
			if len(code) == 0 {
				continue
			}

			// append to the file content after a new line at the end
			content = append(content, []byte("\n"+code)...)

			if err := os.WriteFile(fileLoc, content, 0644); err != nil {
				logger.Error("Error: ", err)
				return err
			}
		}
	}

	return nil
}

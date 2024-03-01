package spawn

import (
	"fmt"
	"log/slog"
)

// FileType tells the application which type of proto file is it so we can sort Txs from Queries
type FileType string

const (
	Tx    FileType = "tx"
	Query FileType = "query"
	None  FileType = "none"
)

// get the str of FileType
func (ft FileType) String() string {
	return string(ft)
}

// ModuleMapping a map of the module name to a list of ProtoRPCs
type ModuleMapping map[string][]*ProtoRPC

func (mm ModuleMapping) Print(logger *slog.Logger) {
	for name, v := range mm {
		v := v
		name := name

		for _, rpc := range v {
			logger.Debug("module", "module", name, "rpc", rpc.Name, "req", rpc.Req, "res", rpc.Res, "module", rpc.Module, "location", rpc.Location, "ftype", rpc.FType, "fileloc", rpc.FileLoc)
		}
	}
}

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

func (pr *ProtoRPC) String() string {
	return fmt.Sprintf(
		"Name: %s, Req: %s, Res: %s, Module: %s, Location: %s, FType: %s, FileLoc: %s",
		pr.Name, pr.Req, pr.Res, pr.Module, pr.Location, pr.FType, pr.FileLoc,
	)
}

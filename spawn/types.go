package spawn

// FileType tells the application which type of proto file is it so we can sort Txs from Queries
type FileType string

const (
	Tx    FileType = "tx"
	Query FileType = "query"
	None  FileType = "none"
)

// ModuleMapping a map of the module name to a list of ProtoRPCs
type ModuleMapping map[string][]*ProtoRPC

func (mm ModuleMapping) Print() {
	for name, v := range mm {
		println(name)
		for _, rpc := range v {
			println(rpc.Name, rpc.FileLoc)
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

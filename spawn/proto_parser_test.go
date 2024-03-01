package spawn

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: false,
	Level:     slog.LevelError,
}))

func TestParser(t *testing.T) {
	type tcase struct {
		name         string
		modPkg       string
		protoContent string
		ft           FileType
		expected     []*ProtoRPC
	}

	tests := []tcase{
		{
			name:   "query, multi line rpc",
			modPkg: "github.com/orgName/chainName",
			ft:     Query,
			protoContent: `syntax = "proto3";
			package cnd.v1;
			import "google/api/annotations.proto";
			import "cnd/v1/genesis.proto";

			option go_package = "github.com/orgName/chainName/x/cnd/types";

			service Query {
				// rpc for Params, make sure this line does not get picked up
				rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
				option (google.api.http).get = "/cnd/v1/params";
				}

				// multi line return statement (from proto linting long names)
				rpc FeeShare(QueryFeeShareRequest)
					returns (QueryFeeShareResponse) {
				option (google.api.http).get = "/cnd/v1/feeshare";
				}
			}`,
			expected: []*ProtoRPC{
				{
					Name:     "Params",
					Req:      "QueryParamsRequest",
					Res:      "QueryParamsResponse",
					Module:   "cnd",
					Location: "x/cnd/types",
					FType:    Query,
					FileLoc:  "",
				},
				{
					Name:     "FeeShare",
					Req:      "QueryFeeShareRequest",
					Res:      "QueryFeeShareResponse",
					Module:   "cnd",
					Location: "x/cnd/types",
					FType:    Query,
					FileLoc:  "",
				},
			},
		},
		{
			name:   "tx, multiple msgs",
			modPkg: "github.com/aaa/bbb",
			ft:     Tx,
			protoContent: `syntax = "proto3";
		package amm.v1;
		import "cosmos/msg/v1/msg.proto";
		import "amm/v1/genesis.proto";
		option go_package = "github.com/aaa/bbb/x/amm/nested/types";

		service Msg {
			option (cosmos.msg.v1.service) = true;

			//rpc does some things
			rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);

			rpc UpdateParams2(MsgUpdateParams2) returns (MsgUpdateParamsResponse2);
		}`,
			expected: []*ProtoRPC{
				{
					Name:     "UpdateParams",
					Req:      "MsgUpdateParams",
					Res:      "MsgUpdateParamsResponse",
					Module:   "amm",
					Location: "x/amm/nested/types",
					FType:    Tx,
					FileLoc:  "",
				},
				{
					Name:     "UpdateParams2",
					Req:      "MsgUpdateParams2",
					Res:      "MsgUpdateParamsResponse2",
					Module:   "amm",
					Location: "x/amm/nested/types",
					FType:    Tx,
					FileLoc:  "",
				},
			},
		},
	}

	defer os.Remove("go.mod")

	for _, tc := range tests {
		tc := tc

		content := []byte(tc.protoContent)

		buildMockGoMod(t, tc.modPkg)

		goPkgDir := GetGoPackageLocationOfFiles(content)

		r := ProtoServiceParser(logger, content, goPkgDir, tc.ft)

		require.Equal(t, len(tc.expected), len(r), tc.name, *r[0])

		require.Equal(t, tc.expected, r)

		os.Remove("go.mod")
	}
}

func TestProtoGoPackageReduction(t *testing.T) {
	type tcase struct {
		modPkg   string
		input    string
		expected string
	}

	tests := []tcase{
		{
			modPkg:   "github.com/rollchains/mychain",
			input:    `option go_package = "github.com/rollchains/mychain/x/cnd/types";`,
			expected: "x/cnd/types",
		},
		{
			modPkg:   "github.com/rollchains/other",
			input:    `option go_package = "github.com/rollchains/other/x/cnd/types";`,
			expected: "x/cnd/types",
		},
		{
			modPkg:   "github.com/rollchains/other",
			input:    `option go_package = "github.com/rollchains/other/x/cosmosmod/types";`,
			expected: "x/cosmosmod/types",
		},
		{
			modPkg:   "github.com/abcchain/misc",
			input:    `option go_package = "github.com/abcchain/misc/x/module/subfolder/types";`,
			expected: "x/module/subfolder/types",
		},
	}

	defer os.Remove("go.mod")

	for _, tc := range tests {
		tc := tc

		buildMockGoMod(t, tc.modPkg)

		actual := GetGoPackageLocationOfFiles([]byte(tc.input))
		require.Equal(t, tc.expected, actual)

		os.Remove("go.mod")
	}
}

func TestBuildProtoInterfaceStub(t *testing.T) {
	type tcase struct {
		pr       ProtoRPC
		expected string
	}

	// take note of the extra line after the final close brace }
	tests := []tcase{
		{
			pr: ProtoRPC{
				Name:   "RPCMethodName",
				Req:    "Query...Request",
				Res:    "Query...Response",
				Module: "mymodule",
				FType:  Query,
			},
			expected: `// RPCMethodName implements types.QueryServer.
func (k Querier) RPCMethodName(goCtx context.Context, req *types.Query...Request) (*types.Query...Response, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	panic("RPCMethodName is unimplemented")
	return &types.Query...Response{}, nil
}
`,
		},
		{
			pr: ProtoRPC{
				Name:   "UpdateParams",
				Req:    "MsgUpdateParams",
				Res:    "MsgUpdateParamsResponse",
				Module: "module",
				FType:  Tx,
			},
			expected: `// UpdateParams implements types.MsgServer.
func (ms msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	panic("UpdateParams is unimplemented")
	return &types.MsgUpdateParamsResponse{}, nil
}
`,
		},
	}

	for _, tc := range tests {
		tc := tc

		res := tc.pr.BuildProtoInterfaceStub()
		require.Equal(t, tc.expected, res)
	}
}

// make sure to `defer os.Remove("go.mod")` after calling
func buildMockGoMod(t *testing.T, moduleName string) {
	// create a go.mod file for this test
	f, err := os.Create("go.mod")
	require.NoError(t, err)
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("module %s", moduleName))
	require.NoError(t, err)
}

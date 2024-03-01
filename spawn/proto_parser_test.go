package spawn

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {

	proto := `syntax = "proto3";
package cnd.v1;

import "google/api/annotations.proto";
import "cnd/v1/genesis.proto";

option go_package = "github.com/rollchains/mychain/x/cnd/types";

// Query provides defines the gRPC querier service.
service Query {
	// Params queries all parameters of the module.
	rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
	option (google.api.http).get = "/cnd/v1/params";
	}

	rpc FeeShare(QueryFeeShareRequest)
	returns (QueryFeeShareResponse) {
	option (google.api.http).get = "/cnd/v1/feeshare";
	}
}`
	fmt.Println(proto)
}

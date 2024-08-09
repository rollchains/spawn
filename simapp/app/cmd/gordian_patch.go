package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	libp2pcryptopb "github.com/libp2p/go-libp2p/core/crypto/pb"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// It appears that libp2p has a bit of a mismatch
// on some protobuf files' names and import paths.
// More specifically, the crypto.proto file is imported as
// the long path "core/crypto/pb/crypto.proto",
// but it is "declared" as the short path "pb/crypto.proto".
// Presumably, most applications importing libp2p don't care about this.
// But the SDK reaches into the global protobuf registry
// and discovers that the mentioned short import was "not declared",
// so early steps (specifically auto CLI setup) fail immediately.
//
// After multiple failed attempts to fix this in any other clean way
// without fixing any external source code,
// here we create a "stub" that does a "public import" of the actual crypto protobufs,
// and then we register it with the global proto registry,
// so that when the SDK's protobuf handling code follows all dependencies
// of all imported files, it can resolve properly.
//
// One of the ugliest hacks I've had to write in recent memory.
func init() {
	const shortPath = "pb/crypto.proto"
	const longPath = "core/crypto/pb/crypto.proto"

	// It seems like protoreflect's builder package would be more suitable
	// than parsing some hardcoded source code,
	// but there does not appear to be an API to create a file descriptor
	// and set a public import,
	// so we have to fall back to parsing source.
	//
	// Note that the public import effectively treats this "file"
	// as a symlink towards the public import.
	parser := protoparse.Parser{
		Accessor: func(filename string) (io.ReadCloser, error) {
			if filename != longPath {
				return nil, fmt.Errorf("unexpected accessed filename %q", filename)
			}

			return io.NopCloser(strings.NewReader(`
syntax = "proto3";

package libp2p_crypto_stub.pb;

import public "pb/crypto.proto";
`)), nil
		},
		LookupImport: func(path string) (*desc.FileDescriptor, error) {
			if path != shortPath {
				return nil, fmt.Errorf("unexpected import path %q", path)
			}

			return desc.WrapFile(libp2pcryptopb.File_pb_crypto_proto)
		},
	}

	fds, err := parser.ParseFiles(longPath)
	if err != nil {
		panic(fmt.Errorf("failed to parse stub file: %w", err))
	}
	if len(fds) != 1 {
		panic(fmt.Errorf("expected 1 parsed file but got %d", len(fds)))
	}

	protoregistry.GlobalFiles.RegisterFile(fds[0].UnwrapFile())
}

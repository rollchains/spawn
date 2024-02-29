package spawn

import (
	"fmt"
	"strings"
)

// The Name, RequestMsg, and ResponseMsg fields of rpc services
type ProtoService struct {
	Name string
	Req  string
	Res  string
}

// ProtoServiceParser parses out a proto file and returns all the services within it.
func ProtoServiceParser(content []byte) []ProtoService {
	qss := make([]ProtoService, 0)
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
			qss = append(qss, ProtoService{
				Name: words[0],
				Req:  words[1],
				Res:  words[2],
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

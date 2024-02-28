go build -buildmode=plugin -o plugins/example.so example-plugin/example.go
chmod +xrw plugins/example.so
EXPORT_LOC=$HOME/.spawn/plugins
mkdir -p $EXPORT_LOC

go build -buildmode=plugin -o $EXPORT_LOC/example.so plugins/example/example-plugin.go
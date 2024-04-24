# sh plugins/example/build.sh

EXPORT_LOC=$HOME/.spawn/plugins
mkdir -p $EXPORT_LOC

NAME="example-plugin"

go build -gcflags="all=-N -l" -mod=readonly -trimpath -o $EXPORT_LOC/$NAME plugins/example/$NAME.go
echo "Plugin built and exported to $EXPORT_LOC/$NAME"
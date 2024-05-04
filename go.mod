module github.com/rollchains/spawn

go 1.21.9

replace github.com/rollchains/spawn/simapp => ./simapp

require (
	github.com/cosmos/btcutil v1.0.5
	github.com/lmittmann/tint v1.0.4
	github.com/manifoldco/promptui v0.9.0
	github.com/mattn/go-isatty v0.0.20
	github.com/rollchains/spawn/simapp v0.0.0-00000000-000000000000
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.9.0
	golang.org/x/mod v0.16.0
	golang.org/x/text v0.14.0
	golang.org/x/tools v0.19.0
)

require (
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/sys v0.19.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

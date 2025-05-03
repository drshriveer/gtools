module github.com/drshriveer/gtools/gomono

go 1.23.0

toolchain go1.23.7

require (
	github.com/drshriveer/gtools/gsync v0.0.0-20250417202014-260178ef6ec0
	github.com/drshriveer/gtools/set v0.0.0-20250422183054-c8fb607d0a7a
	github.com/jessevdk/go-flags v1.6.1
	github.com/stretchr/testify v1.9.0
	golang.org/x/mod v0.24.0
)

replace github.com/drshriveer/gtools/gsync => ../gsync

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

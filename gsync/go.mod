module github.com/drshriveer/gtools/gsync

go 1.23.0

toolchain go1.23.7

require (
	github.com/drshriveer/gtools/gerror v0.0.0
	github.com/stretchr/testify v1.9.0
)

replace github.com/drshriveer/gtools/gerror v0.0.0 => ../gerror

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

module github.com/drshriveer/gtools/gencommon

go 1.23

require (
	github.com/drshriveer/gtools/set v0.0.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/tools v0.19.0
)

replace github.com/drshriveer/gtools/set => ../set

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

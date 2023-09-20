module github.com/drshriveer/gtools/gconfig

go 1.21.1

require (
	github.com/drshriveer/gtools/genum v0.0.0
	github.com/drshriveer/gtools/gerrors v0.0.0
	github.com/drshriveer/gtools/rutils v0.0.0
	github.com/drshriveer/gtools/set v0.0.0
	github.com/puzpuzpuz/xsync/v2 v2.5.0
	github.com/stretchr/testify v1.8.4
	gopkg.in/yaml.v3 v3.0.1
)

replace (
	github.com/drshriveer/gtools/genum v0.0.0 => ../genum
	github.com/drshriveer/gtools/gerrors v0.0.0 => ../gerrors
	github.com/drshriveer/gtools/rutils v0.0.0 => ../rutils
	github.com/drshriveer/gtools/set v0.0.0 => ../set
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

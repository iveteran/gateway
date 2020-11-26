module matrix.works/fmx-gateway

go 1.14

replace matrix.works/fmx-common => ../fmx-common

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	matrix.works/fmx-common v0.0.0-00010101000000-000000000000
)

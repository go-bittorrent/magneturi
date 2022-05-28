magneturi is a library for parsing magnet URI into go struct.

[![Go Reference](https://pkg.go.dev/badge/github.com/go-bittorrent/magneturi.svg)](https://pkg.go.dev/github.com/go-bittorrent/magneturi)
[![x](https://github.com/go-bittorrent/magneturi/actions/workflows/ci.yml/badge.svg)](https://github.com/go-bittorrent/magneturi/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/go-bittorrent/magneturi/branch/main/graph/badge.svg?token=S0F73U8F6H)](https://codecov.io/gh/go-bittorrent/magneturi)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-bittorrent/magneturi)](https://goreportcard.com/report/github.com/go-bittorrent/magneturi)

----

## Installation
```bash
go get github.com/go-bittorrent/magneturi
```

## Example

```go
package main

import (
	"fmt"

	"github.com/go-bittorrent/magneturi"
)

func main() {
	parsed, err := magneturi.Parse("magnet:?xt=urn:btih:9b4c1489bfccd8205d152345f7a8aad52d9a1f57&dn=archlinux-2022.05.01-x86_64.iso")
	if err != nil {
		panic(err)
	}

	fmt.Println(parsed.Encoded()) // magnet:?dn=archlinux-2022.05.01-x86_64.iso&xt=urn:btih:9b4c1489bfccd8205d152345f7a8aad52d9a1f57
}
```

## License
MIT
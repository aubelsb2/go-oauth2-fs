# Direct FS (json) storage for [OAuth 2.0](https://github.com/go-oauth2/oauth2)

This is a simple backend for storage for the OAuth 2.0 it was created for testing purposes and is not intended for
publicly facing use of any sort. It is designed to make the contents easily inspect-able for the user and/or tests
while keeping the requirements minimal. If you can I recommend using the internal and provided BuntsDB.

## Usage

``` go
package main

import (
	"github.com/aubelsb2/go-oauth2-fs"
	"github.com/go-oauth2/oauth2/v4/manage"
	"os"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(go_oauth2_fs.New(os.DirFS("./data/")))

}
```


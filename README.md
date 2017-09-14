# ghttp

This is basically a helper for HTTP clients in Go.

You can use common HTTP methods (GET, POST...), and you can pass a payload which will be sent as url-encoded or JSON.

```go
package main

import (
        "fmt"
        "github.com/gradew/ghttp"
)

func main() {
        ghttp.SetInsecure(true)
        payload:=make(map[string]string)
        payload["username"]="mike"
        payload["secretkey"]="notmyrealpassword"

        response, body:=ghttp.Do("GET", "http://host", payload, true)
        fmt.Printf("Status: %d\n", response)
        fmt.Printf("Body: %s\n", body)
}
```

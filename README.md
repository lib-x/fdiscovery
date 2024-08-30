# fdiscovery
file based tiny discovery lib
## Usage
#### server side
```go
package main

import (
    "log"
    "time"
    "github.com/lib-x/fdiscovery"
)

func main() {
    disc := discovery.NewFileSystemDiscovery("/tmp/services")
    
    serviceName := "myservice"
    serviceAddress := "localhost:8080"

    err := disc.Register(serviceName, serviceAddress)
    if err != nil {
        log.Fatalf("Failed to register service: %v", err)
    }

    // send heartbeat every 20 seconds
    go func() {
        for {
            time.Sleep(20 * time.Second)
            err := disc.Heartbeat(serviceName)
            if err != nil {
                log.Printf("Failed to send heartbeat: %v", err)
            }
        }
    }()

    // your logic here
    select {}
}
```
#### client side

```go   
package main

import (
    "fmt"
    "log"
    "github.com/lib-x/fdiscovery"
)

func main() {
    disc := discovery.NewFileSystemDiscovery("/tmp/services")

    service, err := disc.Discover("myservice")
    if err != nil {
        log.Fatalf("Failed to discover service: %v", err)
    }

    fmt.Printf("Found service: %s at %s\n", service.Name, service.Address)

    // use the service
}
```   
Client for Xiaomi IoT binary protocol. Project based on https://github.com/dovchinnikov/go-miio.

## Install
```sh
go get github.com/ofen/miio-go
```

## Example:

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/ofen/miio-go"
)

func main() {
    client := miio.New("192.168.0.3:54321")
    client.SetToken("c91034a067f36f4558624e65a6f927a7") // will try to use token from handshake if not set

    // https://home.miot-spec.com/spec/mmgg.pet_waterer.s1
    payload := []map[string]interface{}{
        {"did": "cotton_left_time", "siid": 5, "piid": 1},
        {"did": "fault", "siid": 2, "piid": 1},
        {"did": "filter_left_time", "siid": 3, "piid": 1},
        {"did": "indicator_light", "siid": 4, "piid": 1},
        {"did": "lid_up_flag", "siid": 7, "piid": 4},
        {"did": "location", "siid": 9, "piid": 2},
        {"did": "mode", "siid": 2, "piid": 3},
        {"did": "no_water_flag", "siid": 7, "piid": 1},
        {"did": "no_water_time", "siid": 7, "piid": 2},
        {"did": "on", "siid": 2, "piid": 2},
        {"did": "pump_block_flag", "siid": 7, "piid": 3},
        {"did": "remain_clean_time", "siid": 6, "piid": 1},
        {"did": "timezone", "siid": 9, "piid": 1},
    }

    resp, err := client.GetProperties(payload)
    if err != nil {
        panic(client.Token(), err)
    }

    fmt.Printf("%s\n", resp)
}

```
## Additional resources
* https://github.com/OpenMiHome/mihome-binary-protocol
* https://github.com/rytilahti/python-miio
* https://home.miot-spec.com/

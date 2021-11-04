Golang client and Xiaomi IOT binary protocol implementation

This project based on [@dovchinnikov](https://github.com/dovchinnikov) work: https://github.com/dovchinnikov/go-miio

## Install
```go
go get github.com/ofen/miio-go
```

## Example:

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/ofen/miio-go"
)

type response struct {
    ID            int      `json:"id"`
    Result        []result `json:"result"`
    ExecutionTime int      `json:"exe_time"`
}

type result struct {
    DID   string      `json:"did"`
    SIID  int         `json:"siid"`
    PIID  int         `json:"piid"`
    Code  int         `json:"code"`
    Value interface{} `json:"value"`
}

func main() {
    client := miio.New("192.168.0.3:54321", "a0b1c2d3e4f5a0b1c2d3e4f5a0b1c2d3")
    defer client.Close()

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
        panic(err)
    }

    v := response{}

    if err := json.Unmarshal(resp, &v); err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", v)

}
```
## Additional resources
* https://github.com/OpenMiHome/mihome-binary-protocol
* https://github.com/rytilahti/python-miio
* https://home.miot-spec.com/

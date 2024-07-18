package main

import "fmt"
import "encoding/json"

type User struct {
    Name string      `json:"name"`
    Dob string       `json:"dob"`
    Age int          `json:"age"`
}

func main() {
  u := User{}
  err := json.Unmarshal([]byte(`{"name":"John","dob":"2000-01-01","age":20}`), &u)
  if err != nil {
    panic(err)
  }
  fmt.Printf("%#v\n", u)
}

//type Request struct {
//    Operation string      `json:"operation"`
//    Key string            `json:"key"`
//    Value string          `json:"value"`
//}
//
//func main() {
//    s := `{"operation": "get", "key": "example"}`
//    data := Request{}
//    json.Unmarshal([]byte(s), &data)
//    fmt.Printf("%#v\n", data)
//}

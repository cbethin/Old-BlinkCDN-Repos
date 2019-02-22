package main

import (
  "net/http"
  // "strings"
  "fmt"
  "strconv"
)

func handleHttpResponse(w http.ResponseWriter, r *http.Request) {
  var latency = [3]float64{0.454, 3.433, 9.018}

  if (r.Method == "GET") {
    // query := r.URL.Query()
    // destination := query["dest"][0]
    // source := query["source"][0]
    fmt.Println(r.URL.Path)
    if (r.URL.Path == "getpaths") {
      message := ""
      for i := 0; i < len(latency); i++ {
        message += strconv.FormatFloat(latency[i], 'f', -1, 32)
        message += " "
      }

      w.Write([]byte(message))
    } else if (r.URL.Path == "setpaths") {

    }
  }



}

func main() {
  http.HandleFunc("/", handleHttpResponse)
  if err := http.ListenAndServe(":8081", nil); err != nil {
    panic(err)
  }
}

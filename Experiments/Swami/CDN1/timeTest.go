package main
import (
  "os"
  "log"
  "time"
  "fmt"
)



func main()  {
  timeStamp := time.Now().String()
  fmt.Println(timeStamp)
  f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
      log.Fatal(err)
    }
    if _, err := f.Write([]byte(timeStamp + "\n")); err != nil {
      log.Fatal(err)
    }
    if err := f.Close(); err != nil {
      log.Fatal(err)
    }
}

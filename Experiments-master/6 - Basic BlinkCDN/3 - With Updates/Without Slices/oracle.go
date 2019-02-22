package main

import (
  "fmt"
  "os"
  "./blink"
)

type link struct {
	addr            string
	latency         float64
	loss            int
}

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  blink.StartOracle(os.Args[1])
}

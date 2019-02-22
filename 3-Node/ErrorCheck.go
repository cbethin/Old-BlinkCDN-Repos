package main

import(
  "fmt"
  "os"
  "bufio"
  "strings"
  "strconv"
)

func main()  {
  // Open Test File
  file, err := os.Open("ourfile.txt")
  if err != nil {
    fmt.Println("Could not read file")
  }
  // Scan Files
  scannerTest := bufio.NewScanner(file)

  // Variables for file scan and error Calculating
  var tScan string = ""
  var counter float64= 0
  var errScan float64 = 0
  var errorPerc float64 = 0

  for scannerTest.Scan() {
      tScan = string(scannerTest.Text())
      tScanL := strings.Fields(tScan)
      tScanT, _ := strconv.ParseFloat(tScanL[2],64)
      if counter != tScanT {
        errScan += 1
      }
      counter += 1
    }
    // Calculating error Percent
    errorPerc = (errScan / counter) * 100
    // Prints Percent Error
    fmt.Println(errorPerc)
  }

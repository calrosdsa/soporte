package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    num := rand.Intn(9000) + 1000
    fmt.Println(num)
}
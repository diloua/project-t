package main

import (
    "project-t/router"
)

func main() {
    e := route.Init()
    e.Logger.Fatal(e.Start(":8080"))
}

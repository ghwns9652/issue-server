package main

import "issue-server/routes"

func main() {
    r := routes.SetupRouter()
    r.Run(":8080")
}

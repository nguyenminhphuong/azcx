package main

import "github.com/nguyenminhphuong/azcx/cmd"

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}

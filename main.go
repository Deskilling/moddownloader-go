package main

import (
	"moddownloader/filesystem"
	"moddownloader/util"
)

func Init() {
	util.InitLogger(-4)
}

func main() {
	Init()
	filesystem.CreatePath("abc")
}

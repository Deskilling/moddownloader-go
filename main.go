package main

import (
	"moddownloader/util"
)

func Init() {
	util.InitLogger(-4)
	util.ReadConfig()
}

func main() {
	Init()
}

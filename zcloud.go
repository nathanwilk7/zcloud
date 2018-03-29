package main

import (
	"github.com/nathanwilk7/zcloud/cmd"
	"github.com/nathanwilk7/zcloud/controller"
	"github.com/nathanwilk7/zcloud/out"
)

func main () {
	//cmd.RootCmd.Execute()
	pp := cmd.GetProvParamsFromEnv()
	sp := controller.SyncParams{
		Src: "cloud://zcloud-testing/",
		Dest: "testdata",
	}
	controller.Sync(pp, sp, out.New())
}

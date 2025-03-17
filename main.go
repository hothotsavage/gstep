package main

import (
	"github.com/hothotsavage/gstep/config"
	"github.com/hothotsavage/gstep/job"
	"github.com/hothotsavage/gstep/nacos"
	"github.com/hothotsavage/gstep/route"
	"github.com/hothotsavage/gstep/util/db/DbUtil"
)

func main() {
	config.Setup()
	DbUtil.Setup()
	nacos.Setup()
	route.Setup()
	job.Setup()
}

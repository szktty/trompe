package trompe

import (
	"math/rand"
	"time"
)

func Init() {
	now := time.Now()
	rand.Seed(now.Unix())

	InstallPrims()
	InstallModules()
}

func InstallPrims() {
	SetPrim("show", LibCorePrimShow, 1)
}

func InstallModules() {
	RootModule = NewModule(nil, "")
	InstallLibCore()
}

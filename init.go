package trompe

import (
	"math/rand"
	"time"
)

func Init() {
	now := time.Now()
	rand.Seed(now.Unix())
	InstallModules()
}

func InstallModules() {
	RootModule = NewModule(nil, "")
	InstallLibCore()
}

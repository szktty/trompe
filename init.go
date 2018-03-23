package trompe

func Init() {
	InstallPrims()
	InstallModules()
}

func InstallPrims() {
	SetPrim("show", LibCorePrimShow, 1)
}

func InstallModules() {
	RootModule = NewModule(nil, "", nil)
	InstallLibCore()
}

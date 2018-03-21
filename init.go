package trompe

func Init() {
	InstallPrims()
}

func InstallPrims() {
	SetPrim("show", LibCorePrimShow, 1)
}

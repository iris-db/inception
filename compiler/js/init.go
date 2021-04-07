package js

// initProject initializes a new Javascript project in NodeJS.
func initProject() {
	rootCtl := NewFileCtl("sigma-dist")
	rootCtl.InitDir()
	rootCtl.DispatchCommand("npm", ArgsOption("init", "-y"))
}

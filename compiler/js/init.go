package js

import "encoding/json"

// initProject initializes a new Javascript project in NodeJS.
func initProject() {
	ctl := NewFileCtl("sigma-dist")
	ctl.InitDir()

	pkgJson := map[string]interface{}{
		"name":            "my-sigma-api",
		"license":         "MIT",
		"dependencies":    struct{}{},
		"devDependencies": struct{}{},
	}

	b, err := json.MarshalIndent(pkgJson, "", "  ")
	if err != nil {
		panic(err)
	}

	ctl.WriteToFile("package.json", b)

	ctl.WriteToFile("main.js", nil)
}

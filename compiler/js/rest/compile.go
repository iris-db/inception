package rest

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/web-foundation/sigma-production/api"
	"github.com/web-foundation/sigma-production/compiler"
)

var (
	packages     = []string{"typescript@3.7.2"} // Required packages for any API.
	restPackages = []string{"express@4.17.1"}   // REST API dependencies.
)

// APIOpts are universal API options that are translatable to any API compiler
// implementation.
type APIOpts struct {
	Name   string
	Prefix string
	Type   string
	Port   string
	Models api.ModelSet
}

// CompilationOpts are the compiler options to configure the formatting of the
// Javascript project.
type CompilationOpts struct {
	APIOpts
}

// compilationStage is a stage of Javascript compilation. Compilation stages
// are based off of the directories that are currently being written to. For
// example, a compilation stage could involve writing the routes to the routes
// directory.
type compilationStage struct {
	Name string                                            // Name is the compilation stage name
	Impl func(opts CompilationOpts, ctl *compiler.FileCtl) // Impl is the compilation implementation.
	Ctl  *compiler.FileCtl                                 // Ctl is the FileCtl that is used to perform actions on the directory.
}

// CompileAPI begins the compilation processes of a Javascript API by working
// through compilation stages. The full Javascript project is output to a ready
// to use distributable folder with UNIX and Windows scripts to start the API.
func CompileAPI(opts CompilationOpts) {
	// File controllers.
	initProjectCtl := compiler.NewFileCtl(opts.Name)
	createRoutesCtl := compiler.NewFileCtl("src/routes", initProjectCtl)
	createModelsCtl := compiler.NewFileCtl("src/models", initProjectCtl)

	// Initialize compilation stages.
	stages := []compilationStage{
		{Name: "Init Project", Impl: initProject, Ctl: initProjectCtl},
		{Name: "Create Models", Impl: createModels, Ctl: createModelsCtl},
		{Name: "Create Routes", Impl: createRoutes, Ctl: createRoutesCtl},
	}

	// Execute compilation stages in order.
	for _, s := range stages {
		s.Impl(opts, s.Ctl)
	}
}

//go:embed templates/main.txt
var mainTemplate string

// tsConfig are the required tsconfig.json settings for the API to compile
// the api to Javascript. These defaults must not be overridden as it can cause
// problems with the libraries that are used.
var tsConfig = map[string]interface{}{
	"compilerOptions": map[string]interface{}{
		"target":          "es6",
		"module":          "commonjs",
		"esModuleInterop": true,
		"skipLibCheck":    true,
	},
}

// initProject initializes an npm project and installs the required
// dependencies. Creates a main.ts and tsconfig.json file.
func initProject(opts CompilationOpts, ctl *compiler.FileCtl) {
	// Initialize project with npm and install required dependencies.
	ctl.DispatchCommand("npm", compiler.ArgsOption("init", "-y"))
	ctl.DispatchCommand("npm", compiler.ArgOption("i"), compiler.ArgsOption(restPackages...), compiler.ArgsOption(packages...))

	// Create tsconfig.json and write defaults.
	b, err := json.MarshalIndent(tsConfig, "", "\t")
	if err != nil {
		panic(err)
	}
	ctl.WriteToFile("tsconfig.json", b)

	// Parse the main.ts file template and write it.
	r := compiler.ParseTemplate(mainTemplate, compiler.TemplateValues{
		"API_PORT":   compiler.StrPtr(opts.Port),
		"API_PREFIX": compiler.StrPtr(opts.Prefix),
	})
	ctl.WriteToFile("src/main.ts", []byte(r))
}

// createModels is in charge of the models directory of the project.
// It creates controllers for api models for usage in express routes.
// The controllers are responsible for interacting with the database that was
// chosen to be used.
func createModels(opts CompilationOpts, ctl *compiler.FileCtl) {
	for _, m := range opts.Models {
		println(m.Name)
	}
}

//go:embed templates/router.txt
var routerTemplate string

// createRoutes is in charge of creating all crud express routes.
// The routes are located in the routes folder, with the index.ts file
// exporting the master router. The api prefix will add a base route for the
// master router.
func createRoutes(opts CompilationOpts, ctl *compiler.FileCtl) {
	var r []string

	for _, m := range opts.Models {
		tg := CrudGenerator{Model: m}
		r = append(r, tg.Get())
	}

	var sep string
	if len(opts.Models) > 1 {
		sep = "\n"
	}

	t := compiler.ParseTemplate(routerTemplate, compiler.TemplateValues{
		"ROOT_ROUTES": compiler.StrPtr(strings.Join(r, sep)),
	})
	ctl.WriteToFile("index.ts", []byte(t))
}

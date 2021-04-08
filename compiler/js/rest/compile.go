package rest

import (
	"github.com/web-foundation/sigma-production/api"
	"github.com/web-foundation/sigma-production/compiler"
	"github.com/web-foundation/sigma-production/compiler/js/rest/templates"
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

	// Initialize compilation stages.
	stages := []compilationStage{
		{Name: "Init Project", Impl: initProject, Ctl: initProjectCtl},
		{Name: "Create Routes", Impl: createRoutes, Ctl: createRoutesCtl},
	}

	// Execute compilation stages in order.
	for _, s := range stages {
		s.Impl(opts, s.Ctl)
	}
}

func initProject(opts CompilationOpts, ctl *compiler.FileCtl) {
	ctl.DispatchCommand("npm", compiler.ArgsOption("init", "-y"))
	ctl.DispatchCommand("npm", compiler.ArgOption("i"), compiler.ArgsOption(restPackages...), compiler.ArgsOption(packages...))
	t := compiler.ParseTemplate(templates.Main, compiler.TemplateValues{
		"API_PORT":   compiler.StrPtr(opts.Port),
		"API_PREFIX": compiler.StrPtr(opts.Prefix),
	})
	ctl.WriteToFile("src/main.js", []byte(t))
}

func createRoutes(opts CompilationOpts, ctl *compiler.FileCtl) {
	t := compiler.ParseTemplate(templates.Router, compiler.TemplateValues{
		"ROOT_ROUTES": compiler.StrPtr(""),
	})
	ctl.WriteToFile("index.js", []byte(t))
}

package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DispatchCommandOption is an option for DispatchCommand.
type DispatchCommandOption func(c *dispatchCommandConfig)

// dispatchCommandConfig the config struct for DispatchCommand.
// It is configured through a DispatchCommandOption.
type dispatchCommandConfig struct {
	Name             string
	Args             []string
	WorkingDirectory string // WorkingDirectory is the directory to execute the command relative to the FileCtl directory.
}

// ArgsOption is multiple args for a shell command.
func ArgsOption(args ...string) DispatchCommandOption {
	return func(c *dispatchCommandConfig) {
		c.Args = append(c.Args, args...)
	}
}

// ArgOption is a single arg for a shell command.
func ArgOption(arg string) DispatchCommandOption {
	return func(c *dispatchCommandConfig) {
		c.Args = append(c.Args, arg)
	}
}

// WorkingDirectoryOption sets the working directory for a shell command.
func WorkingDirectoryOption(dir string) DispatchCommandOption {
	return func(c *dispatchCommandConfig) {
		c.WorkingDirectory = dir
	}
}

// DispatchCommand executes a command in the current shell at the specified at
// the specified working directory.
func (f FileCtl) DispatchCommand(name string, opts ...DispatchCommandOption) {
	c := &dispatchCommandConfig{Name: name}
	for _, opt := range opts {
		opt(c)
	}

	if err := os.Chdir(f.concatToPath(c.WorkingDirectory)); err != nil {
		panic(err)
	}

	if err := exec.Command(name, c.Args...).Run(); err != nil {
		panic(err)
	}

	for n := 0; n < strings.Count(c.WorkingDirectory, "/")+1; n++ {
		if err := os.Chdir(".."); err != nil {
			panic(err)
		}
	}
}

// WriteToFile writes a byte slice to a file.
func (f FileCtl) WriteToFile(name string, contents []byte) {
	p := f.concatToPath(name)

	// If it is a directory attempt to initialize path.
	t := strings.Split(name, "/")
	if len(t) > 0 {
		dir := f.concatToPath(strings.Join(t[:len(t)-1], "/"))
		_, err := os.Stat(f.concatToPath(dir))
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0777); err != nil {
				panic(err)
			}
		}
	}

	// Create and write to file.
	file, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.Write(contents); err != nil {
		panic(err)
	}
}

// InitDir initializes the working directory.
func (f FileCtl) InitDir() {
	if err := os.MkdirAll(f.Directory, 0777); err != nil {
		panic(err)
	}
}

// concatToPath joins a file name to the working directory.
func (f FileCtl) concatToPath(name string) string {
	if name != "" {
		return f.Directory + "/" + name
	}
	return f.Directory
}

// FileCtlOption is an option for configuring a FileCtl.
type FileCtlOption func(*FileCtl)

// FileCtl is a utility struct for writing Javascript files.
type FileCtl struct {
	Directory string
}

// NewFileCtl constructs a new FileCtl from a directory and options.
func NewFileCtl(directory string, joins ...*FileCtl) *FileCtl {
	ctl := &FileCtl{Directory: directory}
	for _, j := range joins {
		ctl.Directory = fmt.Sprintf("%s/%s", j.Directory, ctl.Directory)
	}
	ctl.InitDir()
	return ctl
}

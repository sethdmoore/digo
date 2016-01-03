package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/types"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	//"strings"
	//"os/exec"
)

// package scope so we don't have to pass around
var p types.Plugins

func search_plugin_dir() (plugin_path string, plugins []string, e error) {
	found := false
	var files []os.FileInfo
	var err error

loop:
	// search all directories specified
	for _, path := range globals.PLUGIN_PATHS {
		// try to read contents of each directory
		files, err = ioutil.ReadDir(path)
		if err == nil {
			found = true
			fmt.Printf("Found plugins directory \"%s\"\n", path)
			plugin_path = path
			break loop
		}
	}

	if found {
		// iterate over all files in the plugins directory
		for _, file := range files {
			name := file.Name()
			if name[0:1] == "." || name[0:1] == "_" {
				// skip plugin config files and hidden files
				continue
			} else if file.IsDir() {
				// skip directories
				continue
			}

			//spew.Dump(file)
			plugins = append(plugins, name)
			fmt.Printf("%v\n", name)
		}
	} else {
		e = errors.New("No plugins directory found")
	}
	return
}

func register_plugin(dir string, file string) (p *types.Plugin, err error) {
	// register string is hardcoded, always the first argument
	config, err := Exec(dir, file, []string{"register"})
	err = json.Unmarshal(config, &p)
	if err != nil {
		err = errors.New(fmt.Sprintf("Couldn't run \"%s register\"\n", file))
		fmt.Printf("%s\n", err)
		fmt.Printf("%s\n", config)
		return p, err
	}
	//spew.Dump(config)
	return p, err
}

func Register() (found bool) {
	var plugin_files []string
	var plugin *types.Plugin
	var err error

	p.Directory, plugin_files, err = search_plugin_dir()
	fmt.Printf("%v\n", plugin_files)
	if err != nil {
		found = false
		fmt.Printf("WARN: Problem with plugins directory: %s\n", err)
		return
	} else {
		found = true
	}

	for _, plugin_name := range plugin_files {
		plugin, err = register_plugin(p.Directory, plugin_name)
		if err != nil {
			fmt.Printf("Could not register %s: %s\n", plugin_name, err)
		} else {
			p.Plugins[plugin_name] = plugin
		}
	}

	return found
}

func Exec(dir string, command string, arguments []string) (output []byte, err error) {
	// maybe? this will work on windows?
	path := filepath.FromSlash(dir + "/" + command)

	cmd := exec.Command(path, arguments...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("%s\n", output)
	}
	return output, err
}

func Init() *types.Plugins {
	p.Plugins = make(map[string]*types.Plugin)

	success := Register()
	if !success {
		fmt.Printf("No plugin directory found\n")
		places := ""
		for _, place := range globals.PLUGIN_PATHS {
			places += place + "\n"
		}

		fmt.Printf("Looked in %s\n", places)
	}

	return &p
}

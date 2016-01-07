package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/types"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	//"strings"
	//"os/exec"
)

// package scope so we don't have to pass around
var p types.Plugins
var log *logging.Logger

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
			log.Notice("Found plugins directory \"%s\"\n", path)
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
			//log.Debug("%v\n", name)
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
		log.Errorf("%s\n", err)
		log.Debugf("%s\n", config)
		return p, err
	}
	// default to simple type plugin
	if p.Type == "" {
		p.Type = "simple"
	}

	if p.Type == "simple" {
		log.Debugf("Simple plugin %s registered", p.Name)
	} else if p.Type == "json" {
		log.Debug("JSON plugin %s  registered", p.Name)
	} else {
		log.Warningf("Plugin of unknown type registered: %s", p.Type)
		log.Warning("Valid types: simple, json")
	}
	//spew.Dump(config)
	return p, err
}

func Register() (found bool) {
	var plugin_files []string
	var enabled_plugins []string
	var plugin *types.Plugin
	var err error

	p.Directory, plugin_files, err = search_plugin_dir()
	log.Debug("Potential plugins: %v\n", plugin_files)
	if err != nil {
		found = false
		log.Warning("Problem with plugins directory: %s\n", err)
		return
	} else {
		found = true
	}

	for _, plugin_name := range plugin_files {
		plugin, err = register_plugin(p.Directory, plugin_name)
		if err != nil {
			log.Warningf("Could not register %s: %s\n", plugin_name, err)
		} else {
			p.Plugins[plugin_name] = plugin
			enabled_plugins = append(enabled_plugins, plugin.Name)
		}
	}

	log.Noticef("Enabled Plugins: %s", strings.Join(enabled_plugins, ", "))

	return found
}

func Exec(dir string, command string, arguments []string) (output []byte, err error) {
	// maybe? this will work on windows?
	path := filepath.FromSlash(dir + "/" + command)

	cmd := exec.Command(path, arguments...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Warningf("%s\n", err)
		log.Debugf("%s\n", output)
	}
	return output, err
}

func ExecJson(dir string, command string, arguments *types.PluginMessage) {
	/*
		path := filepath.FromSlash(dir + "/" + command)
		cmd := exec.Command(path, z)
	*/
}

func Init(logger *logging.Logger) *types.Plugins {
	p.Plugins = make(map[string]*types.Plugin)
	log = logger

	success := Register()
	if !success {
		log.Warning("No plugin directory found")
		places := ""
		for _, place := range globals.PLUGIN_PATHS {
			places += place + "\n"
		}

		log.Warningf("Looked in %s\n", places)
	}

	return &p
}

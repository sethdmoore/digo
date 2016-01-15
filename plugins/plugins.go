package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/config"
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

func searchPluginDir() (pluginPath string, plugins []string, e error) {
	found := false
	var files []os.FileInfo
	var err error

loop:
	// search all directories specified
	for _, path := range globals.PluginPaths {
		// try to read contents of each directory
		files, err = ioutil.ReadDir(path)
		if err == nil {
			found = true
			log.Notice("Found plugins directory \"%s\"\n", path)
			pluginPath = path
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

func registerPlugin(dir string, file string) (plugin *types.Plugin, err error) {
	c := config.Get()
	// register string is hardcoded, always the first argument
	config, err := Exec(dir, file, []string{"register"})
	err = json.Unmarshal(config, &plugin)
	if err != nil {
		err = fmt.Errorf("Couldn't run \"%s register\"\n", file)
		log.Errorf("%s\n", err)
		log.Debugf("%s\n", config)
		return
	}
	// default to simple type plugin
	if plugin.Type == "" {
		plugin.Type = "simple"
	}

	if plugin.Filename == "" {
		plugin.Filename = file
	}

	// input validation
	if len(plugin.Triggers) == 0 && len(plugin.Tokens) == 0 {
		err = fmt.Errorf("Plugin \"%s\" does nothing! It has no triggers or tokens. Not registering.", file)
		return
	}

	if plugin.Type == "simple" {
		log.Debugf("Simple plugin %s registered", plugin.Name)
	} else if plugin.Type == "json" {
		log.Debug("JSON plugin %s registered", plugin.Name)
	} else {
		log.Warningf("Plugin of unknown type registered: %s", plugin.Type)
		log.Warning("Valid types: simple, json")
		err = errors.New("Unknown plugin type")
		return
	}

	// trigger cache so we only have to iterate one set per message
	for _, trigger := range plugin.Triggers {
		// sorry, can't override /bot
		if trigger == c.Trigger {
			log.Info("Prevented plugin %s from trying to override bot trigger %s.", plugin.Name, c.Trigger)
			continue
		}
		p.AllTriggers[trigger] = file
	}

	// token cache so we only have to iterate one set per message
	for _, token := range plugin.Tokens {
		p.AllTokens[token] = file
	}
	//spew.Dump(config)
	return
}

// Register function builds the Plugins struct and calls "register" on each
// plugin
func Register() (found bool) {
	var pluginFiles []string
	var enabledPlugins []string
	var plugin *types.Plugin
	var err error
	c := config.Get()

	// build top level trigger cache
	p.AllTriggers = map[string]string{
		c.Trigger: "__internal",
	}

	p.Directory, pluginFiles, err = searchPluginDir()
	log.Debug("Potential plugins: %v\n", pluginFiles)
	if err != nil {
		found = false
		log.Warning("Problem with plugins directory: %s\n", err)
		return
	}
	found = true

	for _, pluginName := range pluginFiles {
		plugin, err = registerPlugin(p.Directory, pluginName)
		if err != nil {
			log.Warningf("Could not register %s: %s\n", pluginName, err)
		} else {
			p.Plugins[pluginName] = plugin
			enabledPlugins = append(enabledPlugins, plugin.Name)
		}
	}

	log.Noticef("Enabled Plugins: %s", strings.Join(enabledPlugins, ", "))

	return found
}

// Exec executes a command on a simple-type plugin
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

// ExecJSON executes a command on json-type plugin
func ExecJSON(dir string, command string, arguments *types.PluginMessage) ([]byte, error) {
	path := filepath.FromSlash(dir + "/" + command)
	var err error
	var blob []byte

	blob, err = json.Marshal(arguments)

	cmd := exec.Command(path, "json", string(blob))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Warningf("%s\n", err)
		log.Debugf("%s\n", output)
		return []byte{}, err
	}
	log.Debugf("%v", string(output))
	return output, err
}

// Init sets up the Plugins struct and logger
func Init(logger *logging.Logger) *types.Plugins {
	//c := config.Get()
	p.Plugins = make(map[string]*types.Plugin)

	log = logger

	success := Register()
	if !success {
		log.Warning("No plugin directory found")
		places := ""
		for _, place := range globals.PluginPaths {
			places += place + "\n"
		}

		log.Warningf("Looked in %s\n", places)
	}

	return &p
}

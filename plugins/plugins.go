package plugins

import (
	"github.com/sethdmoore/digo/types"
)

func Init() *types.Plugins {
	var p types.Plugins
	p.Plugin = make(map[string]*types.Plugin)
	//Plugins := make(map[string][]Plugin)
	return &p
}

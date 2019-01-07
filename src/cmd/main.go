package main

import (
	"github.com/abilioesteves/metrics-generator-tabajara/src/generator"
	"github.com/abilioesteves/metrics-generator-tabajara/src/hook"
	"github.com/abilioesteves/metrics-generator-tabajara/src/metrics"
)

func main() {
	g := generator.NewGeneratorTabajara(metrics.Init(), generator.GetDefaultEntropy())
	go g.Init()     // fire metrics generator
	go hook.Init(g) // fire webhook
	select {}       // keep-alive magic
}

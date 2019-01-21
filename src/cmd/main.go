package main

import (
	"context"

	"github.com/abilioesteves/goh/gohcmd"
	"github.com/abilioesteves/metrics-generator-tabajara/src/generator"
	"github.com/abilioesteves/metrics-generator-tabajara/src/hook"
	"github.com/abilioesteves/metrics-generator-tabajara/src/metrics"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g := generator.NewGeneratorTabajara(metrics.Init(), generator.GetDefaultEntropy())
	go g.Init(ctx)                   // fire metrics generator
	go hook.NewDefaultHook(g).Init() // fire webhook
	go gohcmd.GracefulStop(cancel)
	select {} // keep-alive magic
}

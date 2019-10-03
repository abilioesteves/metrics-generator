package cmd

import (
	"context"

	"github.com/labbsr0x/goh/gohcmd"
	"github.com/abilioesteves/metrics-generator/generator"
	"github.com/abilioesteves/metrics-generator/hook"
	"github.com/abilioesteves/metrics-generator/metrics"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	g := generator.NewGeneratorTabajara(metrics.Init(), generator.GetDefaultEntropy())
	go g.Init(ctx)                   // fire metrics generator
	go hook.NewDefaultHook(g).Init() // fire webhook
	go gohcmd.GracefulStop(cancel)
	select {} // keep-alive magic
}
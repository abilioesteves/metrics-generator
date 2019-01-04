package main

import (
	"github.com/abilioesteves/metrics-generator-tabajara/src/generator"
	"github.com/abilioesteves/metrics-generator-tabajara/src/hook"
	"github.com/abilioesteves/metrics-generator-tabajara/src/metrics"
)

func main() {
	g := generator.New(metrics.Init())
	go g.Init()
	go hook.Init(g)
	select {}
}

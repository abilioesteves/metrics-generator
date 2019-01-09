package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abilioesteves/metrics-generator-tabajara/src/generator"
	"github.com/abilioesteves/metrics-generator-tabajara/src/hook"
	"github.com/abilioesteves/metrics-generator-tabajara/src/metrics"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g := generator.NewGeneratorTabajara(metrics.Init(), generator.GetDefaultEntropy())
	go g.Init(ctx)         // fire metrics generator
	go hook.InitDefault(g) // fire webhook
	go gracefulStop(cancel)
	select {} // keep-alive magic
}

// gracefullStop cancels gracefully the running goRoutines
func gracefulStop(cancel context.CancelFunc) {
	stopCh := make(chan os.Signal)

	signal.Notify(stopCh, syscall.SIGTERM)
	signal.Notify(stopCh, syscall.SIGINT)

	<-stopCh // waits for a stop signal
	stop(0, cancel)
}

// stop stops this program
func stop(returnCode int, cancel context.CancelFunc) {
	logrus.Infof("Stopping generator...")
	cancel()
	time.Sleep(2 * time.Second)
	os.Exit(returnCode)
}

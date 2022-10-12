package mq

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/suosi-inc/go-demo/mq/internal/mq/setup"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
)

var (
	Cfg = config.Cfg
)

// NewApp New app
func NewApp(args []string) error {
	bootstrap()

	setupDi()

	log.Info("args:", log.Any("", args))

	// Run args
	if len(args) > 0 {

		newRun(args)
	}

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return nil
}

// bootstrap Bootstrap app
func bootstrap() {
	// Config map into struct
	err := viper.Unmarshal(&Cfg)
	if err != nil {
		panic("Unable to decode config into struct: ")
	}
	log.Info("Config into struct", log.Any("cfg", Cfg))
}

// setupDi setup service and set di
func setupDi() {
	setup.Rabbit()
}

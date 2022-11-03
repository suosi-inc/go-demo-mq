package mq

import (
	"github.com/spf13/viper"
	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
)

var (
	Cfg = config.Cfg
)

// NewApp New app
func NewApp(args []string) error {
	bootstrap()

	setupDi()

	log.Infof("args: %v", args)

	// Run args
	if len(args) > 0 {
		newRun(args)
	}

	return nil
}

// bootstrap Bootstrap app
func bootstrap() {
	// Config map into struct
	err := viper.Unmarshal(&Cfg)
	if err != nil {
		panic("Unable to decode config into struct: ")
	}
	log.Info("Config into struct: ", log.Any("cfg", Cfg))
}

// setupDi setup service and set di
func setupDi() {
}

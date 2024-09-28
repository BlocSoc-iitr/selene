package main

import (
	"context"
	"flag"
	
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/BlocSoc-iitr/selene/config"
	"github.com/BlocSoc-iitr/selene/client"
	"github.com/BlocSoc-iitr/selene/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	cliConfig := parseFlags()

	setupLogging()

	cfg := getConfig(cliConfig)
	
	cl, err := client.NewClient(cfg)
	if err != nil {
		logrus.Fatalf("Failed to create client: %v", err)
	}

	if err := cl.Start(); err != nil {
		logrus.Fatalf("Failed to start client: %v", err)
	}

	registerShutdownHandler(cl)

	// Block forever
	select {}
}

func ParseFlags() *config.CliConfig {
	cliConfig := &config.CliConfig{}

	flag.StringVar(cliConfig.ExecutionRpc, "execution-rpc", "", "Execution RPC URL")
	flag.StringVar(cliConfig.ConsensusRpc, "consensus-rpc", "", "Consensus RPC URL")
	flag.StringVar(cliConfig.RpcBindIp, "rpc-bind-ip", "", "RPC bind IP address")
	
	var rpcPort uint
	flag.UintVar(&rpcPort, "rpc-port", 0, "RPC port")
	if rpcPort != 0 {
		port := uint16(rpcPort)
		cliConfig.RpcPort = &port
	}

	flag.StringVar(cliConfig.DataDir, "data-dir", "", "Data directory")
	flag.StringVar(cliConfig.Fallback, "fallback", "", "Fallback URL")

	var loadExternalFallback bool
	flag.BoolVar(&loadExternalFallback, "load-external-fallback", false, "Load external fallback")
	cliConfig.LoadExternalFallback = &loadExternalFallback

	var strictCheckpointAge bool
	flag.BoolVar(&strictCheckpointAge, "strict-checkpoint-age", false, "Strict checkpoint age")
	cliConfig.StrictCheckpointAge = &strictCheckpointAge

	var checkpoint string
	flag.StringVar(&checkpoint, "checkpoint", "", "Checkpoint")
	if checkpoint != "" {
		checkpointBytes, err := utils.Hex_str_to_bytes(checkpoint)
		if err != nil {
			logrus.Fatalf("Failed to parse checkpoint: %v", err)
		}
		cliConfig.Checkpoint = &checkpointBytes
	}

	flag.Parse()

	return cliConfig
}

func SetupLogging() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
}

func getConfig(cliConfig *config.CliConfig) *config.Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatalf("Failed to get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".selene", "config.toml")
	network := "mainnet" // You might want to make this configurable

	cfg := config.Config{}.FromFile(&configPath, &network, cliConfig)

	return &cfg
}

func RegisterShutdownHandler(cl *client.Client) {
	var shutdownCounter int
	var mu sync.Mutex

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range c {
			mu.Lock()
			shutdownCounter++
			counter := shutdownCounter
			mu.Unlock()

			if counter == 3 {
				logrus.Info("Forced shutdown")
				os.Exit(0)
			}

			logrus.Infof("Shutting down... press ctrl-c %d more times to force quit", 3-counter)

			if counter == 1 {
				go func() {
					if err := cl.Shutdown(context.Background()); err != nil {
						logrus.Errorf("Error during shutdown: %v", err)
					}
					os.Exit(0)
				}()
			}

			if sig == syscall.SIGTERM {
				if err := cl.Shutdown(context.Background()); err != nil {
					logrus.Errorf("Error during shutdown: %v", err)
				}
				os.Exit(0)
			}
		}
	}()
}
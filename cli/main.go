package cli

import (
	"fmt"
	// "net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"encoding/hex"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/BlocSoc-iitr/selene/config"
)

// type Config struct {
// 	Network              string
// 	RpcBindIP            net.IP
// 	RpcPort              uint16
// 	Checkpoint           string
// 	ExecutionRPC         string
// 	ConsensusRPC         string
// 	DataDir              string
// 	Fallback             string
// 	LoadExternalFallback bool
// 	StrictCheckpointAge  bool
// }

type Client struct {
	// Add necessary fields
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	config, err := getConfig()
	if err != nil {
		sugar.Errorf("Failed to get config: %v", err)
		os.Exit(1)
	}

	client, err := newClient(config)
	if err != nil {
		sugar.Errorf("Failed to create client: %v", err)
		os.Exit(1)
	}

	if err := client.start(); err != nil {
		sugar.Errorf("Failed to start client: %v", err)
		os.Exit(1)
	}

	registerShutdownHandler(client, sugar)

	// Wait indefinitely
	select {}
}

func registerShutdownHandler(client *Client, logger *zap.SugaredLogger) {
	var shutdownCounter int
	var mu sync.Mutex

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range c {
			_ = sig
			mu.Lock()
			shutdownCounter++
			counter := shutdownCounter
			mu.Unlock()

			if counter == 3 {
				logger.Info("Forced shutdown")
				os.Exit(0)
			}

			logger.Infof("Shutting down... press ctrl-c %d more times to force quit", 3-counter)

			if counter == 1 {
				go func() {
					client.shutdown()
					os.Exit(0)
				}()
			}
		}
	}()
}

func getConfig() (*config.Config, error) {
	cli := &Cli{}

	// Set up command-line flags
	pflag.StringVar(cli.ConsensusRpc, "consensus-rpc", "", "Consensus RPC URL")
	pflag.StringVar(cli.ExecutionRpc, "execution-rpc", "", "Execution RPC URL")
	
	var rpcBindIp string
	pflag.StringVar(&rpcBindIp, "rpc-bind-ip", "", "RPC bind IP")
	
	var rpcPort uint16
	pflag.Uint16Var(&rpcPort, "rpc-port", 0, "RPC port")
	
	var checkpointStr string
	pflag.StringVar(&checkpointStr, "checkpoint", "", "Checkpoint (32 byte hex)")
	
	var dataDir string
	pflag.StringVar(&dataDir, "data-dir", "", "Data directory")
	
	var fallback string
	pflag.StringVar(&fallback, "fallback", "", "Fallback URL")
	
	pflag.BoolVar(&cli.LoadExternalFallback, "load-external-fallback", false, "Load external fallback")
	pflag.BoolVar(&cli.StrictCheckpointAge, "strict-checkpoint-age", false, "Strict checkpoint age")
	
	var databaseType string
	pflag.StringVar(&databaseType, "database-type", "", "Database type")

	pflag.Parse()

	// Bind flags to viper
	viper.BindPFlags(pflag.CommandLine)
	
	// Set up environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Set default values
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %v", err)
		}
		dataDir = filepath.Join(home, ".selene")
	}
	cli.DataDir = &dataDir

	// Process the checkpoint
	if checkpointStr != "" {
		checkpointBytes, err := hex.DecodeString(checkpointStr)
		if err != nil {
			return nil, fmt.Errorf("invalid checkpoint hex string: %v", err)
		}
		if len(checkpointBytes) != 32 {
			return nil, fmt.Errorf("checkpoint must be exactly 32 bytes")
		}
		var checkpoint [32]byte
		copy(checkpoint[:], checkpointBytes)
		cli.Checkpoint = &checkpoint
	}

	// Set pointers for optional fields
	if rpcBindIp != "" {
		cli.RpcBindIp = &rpcBindIp
	}
	if rpcPort != 0 {
		cli.RpcPort = &rpcPort
	}
	if fallback != "" {
		cli.Fallback = &fallback
	}

	cliConfig := cli.asCliConfig()
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".selene", "selene.toml")
	var finalConfig config.Config

	finalConfig.FromFile(&configPath, &cli.Network, &cliConfig)
	
	return &finalConfig, nil
}

type Cli struct {
	Network string
	RpcBindIp *string
	RpcPort *uint16
	Checkpoint *[32]byte
	ExecutionRpc *string
	ConsensusRpc *string
	DataDir *string
	Fallback *string
	LoadExternalFallback bool
	StrictCheckpointAge bool
}

func (c *Cli) asCliConfig() config.CliConfig {
	checkpoint := c.Checkpoint[:]
	return config.CliConfig{
		Checkpoint: &checkpoint,
		ExecutionRpc: c.ExecutionRpc,
		ConsensusRpc: c.ConsensusRpc,
		DataDir: c.DataDir,
		RpcBindIp: c.RpcBindIp,
		RpcPort: c.RpcPort,
		Fallback: c.Fallback,
		LoadExternalFallback: &c.LoadExternalFallback,
		StrictCheckpointAge: &c.StrictCheckpointAge,
	}
}

func newClient(config *config.Config) (*Client, error) {
	// Implement client creation logic
	return &Client{}, nil
}

func (c *Client) start() error {
	// Implement client start logic
	return nil
}

func (c *Client) shutdown() {
	// Implement client shutdown logic
}
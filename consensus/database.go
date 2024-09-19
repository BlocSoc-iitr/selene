package consensus
import (
	"errors"
	"github.com/BlocSoc-iitr/selene/config"
	"os"
	"path/filepath"
)
type Database interface {
	New(cfg *config.BaseConfig) (Database, error)
	SaveCheckpoint(checkpoint []byte) error
	LoadCheckpoint() ([]byte, error)
}
type FileDB struct {
	DataDir           string
	defaultCheckpoint [32]byte
}
func (f *FileDB) New(cfg *config.BaseConfig) (Database, error) {
	if cfg.DataDir == nil || *cfg.DataDir == "" {
		return nil, errors.New("data directory is not set in the config")
	}
	return &FileDB{
		DataDir:           *cfg.DataDir,
		defaultCheckpoint: cfg.DefaultCheckpoint,
	}, nil
}
func (f *FileDB) SaveCheckpoint(checkpoint []byte) error {
	err := os.MkdirAll(f.DataDir, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(f.DataDir, "checkpoint"), checkpoint, 0644)
}
func (f *FileDB) LoadCheckpoint() ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(f.DataDir, "checkpoint"))
	if err != nil {
		if os.IsNotExist(err) {
			return f.defaultCheckpoint[:], nil
		}
		return nil, err
	}
	if len(data) == 32 {
		return data, nil
	}
	return f.defaultCheckpoint[:], nil
}
type ConfigDB struct {
	checkpoint [32]byte
}
func (c *ConfigDB) New(cfg *config.BaseConfig) (Database, error) {
	checkpoint := cfg.DefaultCheckpoint
	if cfg.DataDir == nil {
		return nil, errors.New("data directory is not set in the config")
	}
	return &ConfigDB{
		checkpoint: checkpoint,
	}, nil
}
func (c *ConfigDB) SaveCheckpoint(checkpoint []byte) error {
	return nil
}
func (c *ConfigDB) LoadCheckpoint() ([]byte, error) {
	return c.checkpoint[:], nil
}

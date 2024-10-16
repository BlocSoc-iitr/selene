package consensus

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BlocSoc-iitr/selene/config"
)

func TestFileDB_New(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		DataDir:           &tempDir,
		DefaultCheckpoint: [32]byte{1, 2, 3}, 
	}

	db, err := (&FileDB{}).New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fileDB, ok := db.(*FileDB)
	if !ok {
		t.Fatalf("expected FileDB instance, got %T", db)
	}

	if fileDB.DataDir != tempDir {
		t.Errorf("expected DataDir to be %s, got %s", tempDir, fileDB.DataDir)
	}

	if fileDB.defaultCheckpoint != cfg.DefaultCheckpoint {
		t.Errorf("expected defaultCheckpoint to be %v, got %v", cfg.DefaultCheckpoint, fileDB.defaultCheckpoint)
	}
}

func TestFileDB_SaveAndLoadCheckpoint(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		DataDir:           &tempDir,
		DefaultCheckpoint: [32]byte{1, 2, 3},
	}

	db, err := (&FileDB{}).New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fileDB := db.(*FileDB)

	// Test SaveCheckpoint
	err = fileDB.SaveCheckpoint(cfg.DefaultCheckpoint[:])
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test LoadCheckpoint
	loadedCheckpoint, err := fileDB.LoadCheckpoint()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(loadedCheckpoint) != string(cfg.DefaultCheckpoint[:]) {
		t.Errorf("expected loaded checkpoint to be %s, got %s", cfg.DefaultCheckpoint[:], loadedCheckpoint)
	}

	// Test LoadCheckpoint with default checkpoint when file does not exist
	os.Remove(filepath.Join(fileDB.DataDir, "checkpoint"))
	loadedCheckpoint, err = fileDB.LoadCheckpoint()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(loadedCheckpoint) != string(cfg.DefaultCheckpoint[:]) {
		t.Errorf("expected loaded checkpoint to be default %v, got %v", cfg.DefaultCheckpoint, loadedCheckpoint)
	}
}

func TestConfigDB_New(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		DataDir:           &tempDir,
		DefaultCheckpoint: [32]byte{4, 5, 6},
	}

	db, err := (&ConfigDB{}).New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	configDB, ok := db.(*ConfigDB)
	if !ok {
		t.Fatalf("expected ConfigDB instance, got %T", db)
	}

	if configDB.checkpoint != cfg.DefaultCheckpoint {
		t.Errorf("expected checkpoint to be %v, got %v", cfg.DefaultCheckpoint, configDB.checkpoint)
	}
}

func TestConfigDB_SaveAndLoadCheckpoint(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		DataDir:           &tempDir,
		DefaultCheckpoint: [32]byte{7, 8, 9},
	}

	db, err := (&ConfigDB{}).New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	configDB := db.(*ConfigDB)

	// Test LoadCheckpoint
	loadedCheckpoint, err := configDB.LoadCheckpoint()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(loadedCheckpoint) != string(cfg.DefaultCheckpoint[:]) {
		t.Errorf("expected loaded checkpoint to be %v, got %v", cfg.DefaultCheckpoint, loadedCheckpoint)
	}

	err = configDB.SaveCheckpoint([]byte("some_checkpoint_data"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Ensure that the checkpoint remains the same after calling SaveCheckpoint
	loadedCheckpoint, err = configDB.LoadCheckpoint()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(loadedCheckpoint) != string(cfg.DefaultCheckpoint[:]) {
		t.Errorf("expected loaded checkpoint to be %v, got %v", cfg.DefaultCheckpoint, loadedCheckpoint)
	}
}

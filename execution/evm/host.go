package evm

type Host interface {
	// // Returns a reference to the environment.
	// Env() *Env

	// // Returns a mutable reference to the environment.
	// EnvMut() *Env

	// // Load an account.
	// // Returns (is_cold, is_new_account)
	// LoadAccount(address Address) (LoadAccountResult, bool)

	// // Get the block hash of the given block `number`.
	// BlockHash(number uint64) (B256, bool)

	// // Get balance of `address` and if the account is cold.
	// Balance(address Address) (U256, bool, bool)

	// // Get code of `address` and if the account is cold.
	// Code(address Address) (Bytes, bool, bool)

	// // Get code hash of `address` and if the account is cold.
	// CodeHash(address Address) (B256, bool, bool)

	// // Get storage value of `address` at `index` and if the account is cold.
	// SLoad(address Address, index U256) (U256, bool, bool)

	// // Set storage value of account address at index.
	// // Returns (original, present, new, is_cold).
	// SStore(address Address, index U256, value U256) (SStoreResult, bool)

	// // Get the transient storage value of `address` at `index`.
	// TLoad(address Address, index U256) U256

	// // Set the transient storage value of `address` at `index`.
	// TStore(address Address, index U256, value U256)

	// // Emit a log owned by `address` with given `LogData`.
	// Log(log Log [any])

	// // Mark `address` to be deleted, with funds transferred to `target`.
	// SelfDestruct(address Address, target Address) (SelfDestructResult, bool)
}

// SStoreResult represents the result of an `sstore` operation.
type SStoreResult struct {
	// Value of the storage when it is first read
	OriginalValue U256
	// Current value of the storage
	PresentValue U256
	// New value that is set
	NewValue U256
	// Is storage slot loaded from database
	IsCold bool
}

// LoadAccountResult represents the result of loading an account from the journal state.
type LoadAccountResult struct {
	// Is account cold loaded
	IsCold bool
	// Is account empty; if true, the account is not created.
	IsEmpty bool
}

// SelfDestructResult represents the result of a selfdestruct instruction.
type SelfDestructResult struct {
	HadValue            bool
	TargetExists        bool
	IsCold              bool
	PreviouslyDestroyed bool
}
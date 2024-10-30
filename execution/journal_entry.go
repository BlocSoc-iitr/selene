package evm
import(
	"encoding/json"
	"math/big"
	)
type JournalEntryType uint8
const (
	AccountWarmedType JournalEntryType = iota
	AccountDestroyedType
	AccountTouchedType
	BalanceTransferType
	NonceChangeType
	AccountCreatedType
	StorageChangedType
	StorageWarmedType
	TransientStorageChangeType
	CodeChangeType
)

// Creating a struct for JournalEntry containing all possible fields that might be needed
type JournalEntry struct {
	Type         JournalEntryType
	Address      Address
	Target       Address // Used for AccountDestroyed
	WasDestroyed bool    // Used for AccountDestroyed
	HadBalance   U256    // Used for AccountDestroyed
	Balance      U256    // Used for BalanceTransfer
	From         Address // Used for BalanceTransfer
	To           Address // Used for BalanceTransfer
	Key          U256    // Used for Storage operations
	HadValue     U256    // Used for Storage operations
}

func NewAccountWarmedEntry(address Address) *JournalEntry {
	return &JournalEntry{
		Type:    AccountWarmedType,
		Address: address,
	}
}

// NewAccountDestroyedEntry creates a new journal entry for destroying an account
func NewAccountDestroyedEntry(address, target Address, wasDestroyed bool, hadBalance U256) *JournalEntry {
	return &JournalEntry{
		Type:         AccountDestroyedType,
		Address:      address,
		Target:       target,
		WasDestroyed: wasDestroyed,
		HadBalance:   new(big.Int).Set(hadBalance), //to avoid mutating the original value (had balance not written directly)
	}
}

// NewAccountTouchedEntry creates a new journal entry for touching an account
func NewAccountTouchedEntry(address Address) *JournalEntry {
	return &JournalEntry{
		Type:    AccountTouchedType,
		Address: address,
	}
}

// NewBalanceTransferEntry creates a new journal entry for balance transfer
func NewBalanceTransferEntry(from, to Address, balance U256) *JournalEntry {
	return &JournalEntry{
		Type:    BalanceTransferType,
		From:    from,
		To:      to,
		Balance: new(big.Int).Set(balance),
	}
}

// NewNonceChangeEntry creates a new journal entry for nonce change
func NewNonceChangeEntry(address Address) *JournalEntry {
	return &JournalEntry{
		Type:    NonceChangeType,
		Address: address,
	}
}

// NewAccountCreatedEntry creates a new journal entry for account creation
func NewAccountCreatedEntry(address Address) *JournalEntry {
	return &JournalEntry{
		Type:    AccountCreatedType,
		Address: address,
	}
}

// NewStorageChangedEntry creates a new journal entry for storage change
func NewStorageChangedEntry(address Address, key, hadValue U256) *JournalEntry {
	return &JournalEntry{
		Type:     StorageChangedType,
		Address:  address,
		Key:      new(big.Int).Set(key),
		HadValue: new(big.Int).Set(hadValue),
	}
}

// NewStorageWarmedEntry creates a new journal entry for storage warming
func NewStorageWarmedEntry(address Address, key U256) *JournalEntry {
	return &JournalEntry{
		Type:    StorageWarmedType,
		Address: address,
		Key:     new(big.Int).Set(key),
	}
}

// NewTransientStorageChangeEntry creates a new journal entry for transient storage change
func NewTransientStorageChangeEntry(address Address, key, hadValue U256) *JournalEntry {
	return &JournalEntry{
		Type:     TransientStorageChangeType,
		Address:  address,
		Key:      new(big.Int).Set(key),
		HadValue: new(big.Int).Set(hadValue),
	}
}

// NewCodeChangeEntry creates a new journal entry for code change
func NewCodeChangeEntry(address Address) *JournalEntry {
	return &JournalEntry{
		Type:    CodeChangeType,
		Address: address,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (j JournalEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type         JournalEntryType `json:"type"`
		Address      Address          `json:"address"`
		Target       Address          `json:"target,omitempty"`
		WasDestroyed bool             `json:"was_destroyed,omitempty"`
		HadBalance   U256             `json:"had_balance,omitempty"`
		Balance      U256             `json:"balance,omitempty"`
		From         Address          `json:"from,omitempty"`
		To           Address          `json:"to,omitempty"`
		Key          U256             `json:"key,omitempty"`
		HadValue     U256             `json:"had_value,omitempty"`
	}{
		Type:         j.Type,
		Address:      j.Address,
		Target:       j.Target,
		WasDestroyed: j.WasDestroyed,
		HadBalance:   j.HadBalance,
		Balance:      j.Balance,
		From:         j.From,
		To:           j.To,
		Key:          j.Key,
		HadValue:     j.HadValue,
	})
}
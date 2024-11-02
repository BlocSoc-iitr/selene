package evm

import (
	"encoding/json"
	"fmt"
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

// JournalEntry represents a journal entry with various fields.
type JournalEntry struct {
	Type         JournalEntryType `json:"type"`
	Address      Address          `json:"address"`
	Target       Address          `json:"target,omitempty"`        // Used for AccountDestroyed
	WasDestroyed bool             `json:"was_destroyed,omitempty"` // Used for AccountDestroyed
	HadBalance   U256             `json:"had_balance,omitempty"`   // Used for AccountDestroyed
	Balance      U256             `json:"balance,omitempty"`       // Used for BalanceTransfer
	From         Address          `json:"from,omitempty"`          // Used for BalanceTransfer
	To           Address          `json:"to,omitempty"`            // Used for BalanceTransfer
	Key          U256             `json:"key,omitempty"`           // Used for Storage operations
	HadValue     U256             `json:"had_value,omitempty"`     // Used for Storage operations
}

// MarshalJSON implements the json.Marshaler interface for JournalEntry.
func (j JournalEntry) MarshalJSON() ([]byte, error) {
	type Alias JournalEntry // Create an alias to avoid recursion

	// Helper function to convert U256 to hex string
	u256ToHex := func(u U256) string {
		return fmt.Sprintf("0x%s", (*big.Int)(u).Text(16))
	}

	return json.Marshal(&struct {
		Address    string `json:"address"`
		Target     string `json:"target,omitempty"`
		From       string `json:"from,omitempty"`
		To         string `json:"to,omitempty"`
		Key        string `json:"key,omitempty"`
		HadBalance string `json:"had_balance,omitempty"`
		Balance    string `json:"balance,omitempty"`
		HadValue   string `json:"had_value,omitempty"`
		*Alias
	}{
		Address:    "0x" + fmt.Sprintf("%x", j.Address.Addr[:]), // Convert to hex string
		Target:     "0x" + fmt.Sprintf("%x", j.Target.Addr[:]),  // Convert to hex string
		From:       "0x" + fmt.Sprintf("%x", j.From.Addr[:]),    // Convert to hex string
		To:         "0x" + fmt.Sprintf("%x", j.To.Addr[:]),      // Convert to hex string
		Key:        u256ToHex(j.Key),                            // Convert U256 to hex string
		HadBalance: u256ToHex(j.HadBalance),                     // Convert U256 to hex string
		Balance:    u256ToHex(j.Balance),                        // Convert U256 to hex string
		HadValue:   u256ToHex(j.HadValue),                       // Convert U256 to hex string
		Alias:      (*Alias)(&j),                                // Embed the original struct
	})
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

package evm
import (
	"github.com/BlocSoc-iitr/selene/common"
)
type JournaledState struct {
	State EvmState
	TransientStorage TransientStorage
	Logs []Log[LogData] 
	Depth uint
	Journal [][]JournalEntry
	Spec SpecId
	WarmPreloadedAddresses map[Address]struct{}
}
func NewJournalState(spec SpecId,warmPreloadedAddresses map[Address]struct{}) JournaledState {
	return JournaledState{
		State: nil,
		TransientStorage: nil,
		Logs: []Log[LogData] {},
		Depth: 0,
		Journal: [][]JournalEntry{},
		Spec: spec,
		WarmPreloadedAddresses: warmPreloadedAddresses,
	}
}
type JournalCheckpoint struct {
	Log_i     uint
	Journal_i uint
}
func (j JournaledState)setSpecId(spec SpecId) {
	j.Spec = spec
}
type TransientStorage map[Key]U256
type EvmState map[common.Address]Account
type Key struct {
	Account common.Address
	Slot    U256
}
type Log[T any] struct {
	// The address which emitted this log.
	Address Address `json:"address"`
	// The log data.
	Data T `json:"data"`
}
type LogData struct {
	// The indexed topic list.
	Topics []B256 `json:"topics"`
	// The plain data.
	Data Bytes `json:"data"`
}


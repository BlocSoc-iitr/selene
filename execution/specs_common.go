package evm
func (s SpecId) IsEnabledIn(other SpecId) bool {
	return s >= other
}
type Spec interface {
	SpecID() SpecId
	Enabled(specID SpecId) bool
}
type BaseSpec struct {
	id SpecId
}
func (s BaseSpec) SpecID() SpecId {
	return s.id
}

func (s BaseSpec) Enabled(specID SpecId) bool {
	return s.id >= specID
}
func NewSpec(id SpecId) BaseSpec {
	return BaseSpec{id: id}
}
func TryFromUint8(specID uint8) (SpecId, bool) {
	if specID > uint8(PRAGUE_EOF) && specID != uint8(LATEST) {
		return 0, false
	}
	return SpecId(specID), true
}
type OptimismFields struct {
    SourceHash          *B256  `json:"source_hash,omitempty"`
    Mint               *uint64 `json:"mint,omitempty"`
    IsSystemTransaction *bool   `json:"is_system_transaction,omitempty"`
    EnvelopedTx        Bytes   `json:"enveloped_tx,omitempty"`
}

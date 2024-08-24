package types


type BeaconBlock struct{}
type BeaconBlockBody struct{}// 1 method
func (b BeaconBlockBody) def(){}//default

type SignedBlsToExecutionChange struct{}
type BlsToExecutionChange struct{}

type ExecutionPayload struct{}// 1 method
func (exe ExecutionPayload) def(){}// default

type Withdrawwal struct{}
type ProposalSlashing struct{}
type SignedBeaconBlockHeader struct{}
type AttesterSlashing struct{}
type IndexedAttestation struct{}
type Attestation struct {}
type AttestationData struct{}
type Checkpoint struct{}
type SignedVoluntaryExit struct{}
type VoluntaryExit struct{}
type Deposit struct{}
type DepositData struct{}
type Eth1Data struct{}
type Updatec struct{}
type FinalityUpdate struct{}
type OptimisticUpdate struct{}
type Bootstrap struct{}

type Header struct{}
type SyncCommittee struct{}
type SyncAggregate struct{}

type GenericUpdate struct{}
func (g GenericUpdate)from(){

}
package consensus

///NOTE: only these imports are required others are in package already
// uses rpc
// uses config for networks
// uses common for datatypes
import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"os"
	"sync"
	"time"

	"github.com/BlocSoc-iitr/selene/common"
	"github.com/BlocSoc-iitr/selene/config"
	"github.com/BlocSoc-iitr/selene/config/checkpoints"
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/BlocSoc-iitr/selene/consensus/rpc"

	"github.com/BlocSoc-iitr/selene/utils"
	beacon "github.com/ethereum/go-ethereum/beacon/types"
	geth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	bls "github.com/protolambda/bls12-381-util"
)

// Error definitions

var (
	ErrIncorrectRpcNetwork           = errors.New("incorrect RPC network")
	ErrPayloadNotFound               = errors.New("payload not found")
	ErrInvalidHeaderHash             = errors.New("invalid header hash")
	ErrCheckpointTooOld              = errors.New("checkpoint too old")
	ErrBootstrapFetchFailed          = errors.New("could not fetch bootstrap")
	ErrInvalidUpdate                 = errors.New("invalid update")
	ErrInsufficientParticipation     = errors.New("insufficient participation")
	ErrInvalidTimestamp              = errors.New("invalid timestamp")
	ErrInvalidPeriod                 = errors.New("invalid period")
	ErrNotRelevant                   = errors.New("update not relevant")
	ErrInvalidFinalityProof          = errors.New("invalid finality proof")
	ErrInvalidNextSyncCommitteeProof = errors.New("invalid next sync committee proof")
	ErrInvalidSignature              = errors.New("invalid signature")
)

const MAX_REQUEST_LIGHT_CLIENT_UPDATES = 128

type GenericUpdate struct {
	AttestedHeader          consensus_core.Header
	SyncAggregate           consensus_core.SyncAggregate
	SignatureSlot           uint64
	NextSyncCommittee       *consensus_core.SyncCommittee
	NextSyncCommitteeBranch *[]consensus_core.Bytes32
	FinalizedHeader         consensus_core.Header
	FinalityBranch          []consensus_core.Bytes32
}

type ConsensusClient struct {
	BlockRecv          *common.Block
	FinalizedBlockRecv *common.Block
	CheckpointRecv     *[]byte
	genesisTime        uint64
	db                 Database
}

type Inner struct {
	RPC                rpc.ConsensusRpc
	Store              LightClientStore
	lastCheckpoint     *[]byte
	blockSend          chan common.Block
	finalizedBlockSend chan *common.Block
	checkpointSend     chan *[]byte
	Config             *config.Config
}
type LightClientStore struct {
	FinalizedHeader               consensus_core.Header
	CurrentSyncCommitee           consensus_core.SyncCommittee
	NextSyncCommitee              *consensus_core.SyncCommittee
	OptimisticHeader              consensus_core.Header
	PreviousMaxActiveParticipants uint64
	CurrentMaxActiveParticipants  uint64
}

func (con ConsensusClient) New(rpc *string, config config.Config) ConsensusClient {
	blockSend := make(chan common.Block, 256)
	finalizedBlockSend := make(chan *common.Block)
	checkpointSend := make(chan *[]byte)

	db, err := con.db.New(&config)
	if err != nil {
		panic(err)
	}

	var initialCheckpoint [32]byte

	if config.Checkpoint != nil {
		initialNewCheckpoint, errorWhileLoadingCheckpoint := db.LoadCheckpoint()
		copy(initialCheckpoint[:], initialNewCheckpoint)
		if errorWhileLoadingCheckpoint != nil {
			log.Printf("error while loading checkpoint: %v", errorWhileLoadingCheckpoint)
		}
	}
	if initialCheckpoint == [32]byte{} {
		panic("No checkpoint found")
	}
	In := &Inner{}
	inner := In.New(*rpc, blockSend, finalizedBlockSend, checkpointSend, &config)

	go func() {
		err := inner.sync(initialCheckpoint)
		if err != nil {
			if inner.Config.LoadExternalFallback {
				err = sync_all_fallback(inner, inner.Config.Chain.ChainID)
				if err != nil {
					log.Printf("sync failed: %v", err)
					os.Exit(1)
				}
			} else if inner.Config.Fallback != nil {
				err = sync_fallback(inner, inner.Config.Fallback)
				if err != nil {
					log.Printf("sync failed: %v", err)
					os.Exit(1)
				}
			} else {
				log.Printf("sync failed: %v", err)
				os.Exit(1)
			}
		}

		_ = inner.send_blocks()

		for {
			time.Sleep(inner.duration_until_next_update())

			err := inner.advance()
			if err != nil {
				log.Printf("advance error: %v", err)
				continue
			}

			err = inner.send_blocks()
			if err != nil {
				log.Printf("send error: %v", err)
				continue
			}
		}
	}()

	blocksReceived := <-blockSend
	finalizedBlocksReceived := <-finalizedBlockSend
	checkpointsReceived := <-checkpointSend

	return ConsensusClient{
		BlockRecv:          &blocksReceived,
		FinalizedBlockRecv: finalizedBlocksReceived,
		CheckpointRecv:     checkpointsReceived,
		genesisTime:        config.Chain.GenesisTime,
		db:                 db,
	}

}
func (con ConsensusClient) Shutdown() error {
	checkpoint := con.CheckpointRecv
	if checkpoint != nil {
		err := con.db.SaveCheckpoint(*checkpoint)
		if err != nil {
			return err
		}
	}
	return nil
}
func (con ConsensusClient) Expected_current_slot() uint64 {
	now := time.Now().Unix()
	// Assuming SLOT_DURATION is the duration of each slot in seconds
	const SLOT_DURATION uint64 = 12
	return (uint64(now) - con.genesisTime) / SLOT_DURATION
}

func sync_fallback(inner *Inner, fallback *string) error {
	// Create a buffered channel to receive any errors from the goroutine
	errorChan := make(chan error, 1)

	go func() {
		// Attempt to fetch the latest checkpoint from the API
		cf, err := (&checkpoints.CheckpointFallback{}).FetchLatestCheckpointFromApi(*fallback)
		if err != nil {

			errorChan <- err
			return
		}

		if err := inner.sync(cf); err != nil {

			errorChan <- err
			return
		}

		errorChan <- nil
	}()

	return <-errorChan
}

func sync_all_fallback(inner *Inner, chainID uint64) error {
	var n config.Network
	network, err := n.ChainID(chainID)
	if err != nil {
		return err
	}
	errorChan := make(chan error, 1)

	go func() {

		ch := checkpoints.CheckpointFallback{}

		checkpointFallback, errWhileCheckpoint := ch.Build()
		if errWhileCheckpoint != nil {
			errorChan <- errWhileCheckpoint
			return
		}

		chainId := network.Chain.ChainID
		var networkName config.Network
		if chainId == 1 {
			networkName = config.MAINNET
		} else if chainId == 5 {
			networkName = config.GOERLI
		} else if chainId == 11155111 {
			networkName = config.SEPOLIA
		} else {
			errorChan <- errors.New("chain id not recognized")
			return
		}

		// Fetch the latest checkpoint from the network
		checkpoint := checkpointFallback.FetchLatestCheckpoint(networkName)

		// Sync using the inner struct's sync method
		if err := inner.sync(checkpoint); err != nil {
			errorChan <- err
		}

		errorChan <- nil
	}()
	return <-errorChan
}

func (in *Inner) New(rpcURL string, blockSend chan common.Block, finalizedBlockSend chan *common.Block, checkpointSend chan *[]byte, config *config.Config) *Inner {
	rpcClient := rpc.NewConsensusRpc(rpcURL)

	return &Inner{
		RPC:                rpcClient,
		Store:              LightClientStore{},
		lastCheckpoint:     nil, // No checkpoint initially
		blockSend:          blockSend,
		finalizedBlockSend: finalizedBlockSend,
		checkpointSend:     checkpointSend,
		Config:             config,
	}

}
func (in *Inner) Check_rpc() error {
	errorChan := make(chan error, 1)

	go func() {
		chainID, err := in.RPC.ChainId()
		if err != nil {
			errorChan <- err
			return
		}
		if chainID != in.Config.Chain.ChainID {
			errorChan <- ErrIncorrectRpcNetwork
			return
		}
		errorChan <- nil
	}()
	return <-errorChan
}
func (in *Inner) get_execution_payload(slot *uint64) (*consensus_core.ExecutionPayload, error) {
	errorChan := make(chan error, 1)
	blockChan := make(chan consensus_core.BeaconBlock, 1)
	go func() {
		var err error
		block, err := in.RPC.GetBlock(*slot)
		if err != nil {
			errorChan <- err
		}
		errorChan <- nil
		blockChan <- block
	}()

	if err := <-errorChan; err != nil {
		return nil, err
	}

	block := <-blockChan
	Gethblock, err := beacon.BlockFromJSON("capella", block.Hash)
	if err != nil {
		return nil, err
	}
	blockHash := Gethblock.Root()
	latestSlot := in.Store.OptimisticHeader.Slot
	finalizedSlot := in.Store.FinalizedHeader.Slot

	var verifiedBlockHash geth.Hash
	if *slot == latestSlot {
		verifiedBlockHash = toGethHeader(&in.Store.OptimisticHeader).Hash()
	} else if *slot == finalizedSlot {
		verifiedBlockHash = toGethHeader(&in.Store.FinalizedHeader).Hash()
	} else {
		return nil, ErrPayloadNotFound
	}

	// Compare the hashes
	if !bytes.Equal(verifiedBlockHash[:], blockHash.Bytes()) {
		return nil, fmt.Errorf("%w: expected %v but got %v", ErrInvalidHeaderHash, verifiedBlockHash, blockHash)
	}

	payload := block.Body.ExecutionPayload
	return &payload, nil
}

func (in *Inner) Get_payloads(startSlot, endSlot uint64) ([]interface{}, error) {
	var payloads []interface{}

	// Fetch the block at endSlot to get the initial parent hash
	endBlock, err := in.RPC.GetBlock(endSlot)
	if err != nil {
		return nil, err
	}
	endPayload := endBlock.Body.ExecutionPayload
	prevParentHash := endPayload.ParentHash

	// Create a wait group to manage concurrent fetching
	var wg sync.WaitGroup
	payloadsChan := make(chan interface{}, endSlot-startSlot+1)
	errorChan := make(chan error, 1) // Buffer for one error

	// Fetch blocks in parallel
	for slot := endSlot; slot >= startSlot; slot-- {
		wg.Add(1)
		go func(slot uint64) {
			defer wg.Done()
			block, err := in.RPC.GetBlock(slot)
			if err != nil {
				errorChan <- err
				return
			}
			payload := block.Body.ExecutionPayload
			if payload.ParentHash != prevParentHash {
				log.Printf("Error while backfilling blocks: expected block hash %v but got %v", prevParentHash, payload.ParentHash)
				return
			}
			prevParentHash = payload.ParentHash
			payloadsChan <- payload
		}(slot)
	}

	// Close channels after all fetches are complete
	go func() {
		wg.Wait()
		close(payloadsChan)
		close(errorChan)
	}()

	// Collect results and check for errors
	for {
		select {
		case payload, ok := <-payloadsChan:
			if !ok {
				return nil, errors.New("payloads channel closed unexpectedly")
			}
			payloads = append(payloads, payload)
			return payloads, nil
		case err := <-errorChan:
			return nil, err
		}
	}

}
func (in *Inner) advance() error {
	ErrorChan := make(chan error, 1)
	finalityChan := make(chan consensus_core.FinalityUpdate, 1)

	go func() {
		finalityUpdate, err := in.RPC.GetFinalityUpdate()
		if err != nil {
			ErrorChan <- err
			return
		}
		finalityChan <- finalityUpdate
		ErrorChan <- nil
	}()
	if ErrorChan != nil {
		return <-ErrorChan
	}
	finalityUpdate := <-finalityChan
	if err := in.verify_finality_update(&finalityUpdate); err != nil {
		return err
	}
	in.apply_finality_update(&finalityUpdate)

	// Fetch and apply optimistic update
	optimisticUpdate, err := in.RPC.GetOptimisticUpdate()
	if err != nil {
		return err
	}
	if err := in.verify_optimistic_update(&optimisticUpdate); err != nil {
		return err
	}
	in.apply_optimistic_update(&optimisticUpdate)

	// Check for sync committee update if it's not set
	if in.Store.NextSyncCommitee == nil {
		log.Printf("checking for sync committee update")

		currentPeriod := utils.CalcSyncPeriod(in.Store.FinalizedHeader.Slot)
		updates, err := in.RPC.GetUpdates(currentPeriod, 1)
		if err != nil {
			return err
		}

		if len(updates) == 1 {
			update := updates[0]
			if err := in.verify_update(&update); err == nil {
				log.Printf("updating sync committee")
				in.apply_update(&update)
			}
		}
	}

	return nil
}
func (in *Inner) sync(checkpoint [32]byte) error {
	// Reset store and checkpoint
	in.Store = LightClientStore{}
	in.lastCheckpoint = nil

	// Perform bootstrap with the given checkpoint
	in.bootstrap(checkpoint)

	currentPeriod := utils.CalcSyncPeriod(in.Store.FinalizedHeader.Slot)

	errorChan := make(chan error, 1)
	var updates []consensus_core.Update
	var err error
	go func() {
		updates, err = in.RPC.GetUpdates(currentPeriod, MAX_REQUEST_LIGHT_CLIENT_UPDATES)
		if err != nil {
			errorChan <- err
		}

		// Apply updates
		for _, update := range updates {
			if err := in.verify_update(&update); err != nil {

				errorChan <- err
				return
			}
			in.apply_update(&update)
		}

		finalityUpdate, err := in.RPC.GetFinalityUpdate()
		if err != nil {

			errorChan <- err
			return
		}

		if err := in.verify_finality_update(&finalityUpdate); err != nil {
			errorChan <- err
			return
		}
		in.apply_finality_update(&finalityUpdate)

		// Fetch and apply optimistic update

		optimisticUpdate, err := in.RPC.GetOptimisticUpdate()
		if err != nil {
			errorChan <- err
			return
		}
		if err := in.verify_optimistic_update(&optimisticUpdate); err != nil {
			errorChan <- err
			return
		}
		in.apply_optimistic_update(&optimisticUpdate)
		errorChan <- nil
		log.Printf("consensus client in sync with checkpoint: 0x%s", hex.EncodeToString(checkpoint[:]))
	}()

	if err := <-errorChan; err != nil {
		return err
	}
	// Log the success message

	return nil
}
func (in *Inner) send_blocks() error {
	// Get slot from the optimistic header
	slot := in.Store.OptimisticHeader.Slot
	payload, err := in.get_execution_payload(&slot)
	if err != nil {
		return err
	}

	// Get finalized slot from the finalized header
	finalizedSlot := in.Store.FinalizedHeader.Slot
	finalizedPayload, err := in.get_execution_payload(&finalizedSlot)
	if err != nil {
		return err
	}

	// Send payload converted to block over the BlockSend channel
	go func() {
		block, err := PayloadToBlock(payload)
		if err != nil {
			log.Printf("Error converting payload to block: %v", err)
			return
		}
		in.blockSend <- *block
	}()

	go func() {
		block, err := PayloadToBlock(finalizedPayload)
		if err != nil {
			log.Printf("Error converting finalized payload to block: %v", err)
			return
		}
		in.finalizedBlockSend <- block
	}()

	// Send checkpoint over the CheckpointSend channel
	go func() {
		in.checkpointSend <- in.lastCheckpoint
	}()

	return nil
}

// / Gets the duration until the next update
// / Updates are scheduled for 4 seconds into each slot
func (in *Inner) duration_until_next_update() time.Duration {
	currentSlot := in.expected_current_slot()
	nextSlot := currentSlot + 1
	nextSlotTimestamp := nextSlot*12 + in.Config.Chain.GenesisTime

	now := uint64(time.Now().Unix())
	timeToNextSlot := int64(nextSlotTimestamp - now)
	nextUpdate := timeToNextSlot + 4

	return time.Duration(nextUpdate) * time.Second
}
func (in *Inner) bootstrap(checkpoint [32]byte) {
	errorChan := make(chan error, 1)
	bootstrapChan := make(chan consensus_core.Bootstrap, 1)
	go func() {
		bootstrap, errInBootstrap := in.RPC.GetBootstrap(checkpoint)

		if errInBootstrap != nil {
			log.Printf("failed to fetch bootstrap: %v", errInBootstrap)
			errorChan <- errInBootstrap
			return
		}
		bootstrapChan <- bootstrap
		errorChan <- nil
	}()
	if err := <-errorChan; err != nil {
		return
	}
	bootstrap := <-bootstrapChan

	isValid := in.is_valid_checkpoint(bootstrap.Header.Slot)
	if !isValid {
		if in.Config.StrictCheckpointAge {
			panic("checkpoint too old, consider using a more recent checkpoint")
		} else {
			log.Printf("checkpoint too old, consider using a more recent checkpoint")
		}
	}

	verify_bootstrap(checkpoint, bootstrap)
	apply_bootstrap(&in.Store, bootstrap)

}
func verify_bootstrap(checkpoint [32]byte, bootstrap consensus_core.Bootstrap) {
	isCommitteValid := isCurrentCommitteeProofValid(&bootstrap.Header, &bootstrap.CurrentSyncCommittee, bootstrap.CurrentSyncCommitteeBranch)
	if !isCommitteValid {
		log.Println("invalid current sync committee proof")
		return
	}

	headerHash := toGethHeader(&bootstrap.Header).Hash()

	HeaderValid := bytes.Equal(headerHash[:], checkpoint[:])

	if !HeaderValid {
		log.Println("invalid header hash")
		return
	}

}

func apply_bootstrap(store *LightClientStore, bootstrap consensus_core.Bootstrap) {
	store.FinalizedHeader = bootstrap.Header
	store.CurrentSyncCommitee = bootstrap.CurrentSyncCommittee
	store.NextSyncCommitee = nil
	store.OptimisticHeader = bootstrap.Header
	store.PreviousMaxActiveParticipants = 0
	store.CurrentMaxActiveParticipants = 0

}

func (in *Inner) verify_generic_update(update *GenericUpdate, expectedCurrentSlot uint64, store *LightClientStore, genesisRoots []byte, forks consensus_core.Forks) error {
	{
		bits := getBits(update.SyncAggregate.SyncCommitteeBits[:])
		if bits == 0 {
			return ErrInsufficientParticipation
		}

		updateFinalizedSlot := update.FinalizedHeader.Slot
		validTime := expectedCurrentSlot >= update.SignatureSlot &&
			update.SignatureSlot > update.AttestedHeader.Slot &&
			update.AttestedHeader.Slot >= updateFinalizedSlot

		if !validTime {
			return ErrInvalidTimestamp
		}

		storePeriod := utils.CalcSyncPeriod(store.FinalizedHeader.Slot)
		updateSigPeriod := utils.CalcSyncPeriod(update.SignatureSlot)

		var validPeriod bool
		if store.NextSyncCommitee != nil {
			validPeriod = updateSigPeriod == storePeriod || updateSigPeriod == storePeriod+1
		} else {
			validPeriod = updateSigPeriod == storePeriod
		}

		if !validPeriod {
			return ErrInvalidPeriod
		}

		updateAttestedPeriod := utils.CalcSyncPeriod(update.AttestedHeader.Slot)
		updateHasNextCommittee := store.NextSyncCommitee == nil && update.NextSyncCommittee != nil && updateAttestedPeriod == storePeriod

		if update.AttestedHeader.Slot <= store.FinalizedHeader.Slot && !updateHasNextCommittee {
			return ErrNotRelevant
		}

		// Validate finalized header and finality branch
		if update.FinalizedHeader != (consensus_core.Header{}) && update.FinalityBranch != nil {
			if !isFinalityProofValid(&update.AttestedHeader, &update.FinalizedHeader, update.FinalityBranch) {
				return ErrInvalidFinalityProof
			}
		} else if (update.FinalizedHeader != (consensus_core.Header{}) && update.FinalityBranch == nil) || (update.FinalizedHeader == (consensus_core.Header{}) && update.FinalityBranch != nil) {
			return ErrInvalidFinalityProof
		}

		if update.NextSyncCommittee != nil && update.NextSyncCommitteeBranch != nil {
			if !isNextCommitteeProofValid(&update.AttestedHeader, update.NextSyncCommittee, *update.NextSyncCommitteeBranch) {
				return ErrInvalidNextSyncCommitteeProof
			}
		} else if (update.NextSyncCommittee != nil && update.NextSyncCommitteeBranch == nil) || (update.NextSyncCommittee == nil && update.NextSyncCommitteeBranch != nil) {
			return ErrInvalidNextSyncCommitteeProof
		}

		// Set sync committee based on updateSigPeriod
		var syncCommittee *consensus_core.SyncCommittee
		if updateSigPeriod == storePeriod {
			syncCommittee = &in.Store.CurrentSyncCommitee
		} else {
			syncCommittee = in.Store.NextSyncCommitee
		}
		_, err := utils.GetParticipatingKeys(syncCommittee, [64]byte(update.SyncAggregate.SyncCommitteeBits))
		if err != nil {
			return fmt.Errorf("failed to get participating keys: %w", err)
		}

		forkVersion := utils.CalculateForkVersion(&forks, update.SignatureSlot)
		forkDataRoot := utils.ComputeForkDataRoot(forkVersion, consensus_core.Bytes32(genesisRoots))

		if !verifySyncCommitteeSignature(syncCommittee.Pubkeys, &update.AttestedHeader, &update.SyncAggregate, forkDataRoot) {
			return ErrInvalidSignature
		}

		return nil
	}
}
func (in *Inner) verify_update(update *consensus_core.Update) error {
	genUpdate := GenericUpdate{
		AttestedHeader:          update.AttestedHeader,
		SyncAggregate:           update.SyncAggregate,
		SignatureSlot:           update.SignatureSlot,
		NextSyncCommittee:       &update.NextSyncCommittee,
		NextSyncCommitteeBranch: &update.NextSyncCommitteeBranch,
		FinalizedHeader:         update.FinalizedHeader,
		FinalityBranch:          update.FinalityBranch,
	}
	return in.verify_generic_update(&genUpdate, in.expected_current_slot(), &in.Store, in.Config.Chain.GenesisRoot, in.Config.Forks)
}
func (in *Inner) verify_finality_update(update *consensus_core.FinalityUpdate) error {
	if update == nil {
		return ErrInvalidUpdate
	}
	genUpdate := GenericUpdate{
		AttestedHeader:  update.AttestedHeader,
		SyncAggregate:   update.SyncAggregate,
		SignatureSlot:   update.SignatureSlot,
		FinalizedHeader: update.FinalizedHeader,
		FinalityBranch:  update.FinalityBranch,
	}
	return in.verify_generic_update(&genUpdate, in.expected_current_slot(), &in.Store, in.Config.Chain.GenesisRoot, in.Config.Forks)
}
func (in *Inner) verify_optimistic_update(update *consensus_core.OptimisticUpdate) error {
	genUpdate := GenericUpdate{
		AttestedHeader: update.AttestedHeader,
		SyncAggregate:  update.SyncAggregate,
		SignatureSlot:  update.SignatureSlot,
	}
	return in.verify_generic_update(&genUpdate, in.expected_current_slot(), &in.Store, in.Config.Chain.GenesisRoot, in.Config.Forks)
}
func (in *Inner) apply_generic_update(store *LightClientStore, update *GenericUpdate) *[]byte {
	committeeBits := getBits(update.SyncAggregate.SyncCommitteeBits[:])

	// Update max active participants
	if committeeBits > store.CurrentMaxActiveParticipants {
		store.CurrentMaxActiveParticipants = committeeBits
	}

	// Determine if we should update the optimistic header
	shouldUpdateOptimistic := committeeBits > in.safety_threshold() &&
		update.AttestedHeader.Slot > store.OptimisticHeader.Slot

	if shouldUpdateOptimistic {
		store.OptimisticHeader = update.AttestedHeader
	}

	updateAttestedPeriod := utils.CalcSyncPeriod(update.AttestedHeader.Slot)

	updateFinalizedSlot := uint64(0)
	if update.FinalizedHeader != (consensus_core.Header{}) {
		updateFinalizedSlot = update.FinalizedHeader.Slot
	}
	updateFinalizedPeriod := utils.CalcSyncPeriod(updateFinalizedSlot)

	updateHasFinalizedNextCommittee := in.Store.NextSyncCommitee == nil &&
		in.has_sync_update(update) && in.has_finality_update(update) &&
		updateFinalizedPeriod == updateAttestedPeriod

	// Determine if we should apply the update
	hasMajority := committeeBits*3 >= 512*2
	if !hasMajority {
		log.Println("skipping block with low vote count")
	}

	updateIsNewer := updateFinalizedSlot > store.FinalizedHeader.Slot
	goodUpdate := updateIsNewer || updateHasFinalizedNextCommittee

	shouldApplyUpdate := hasMajority && goodUpdate

	// Apply the update if conditions are met
	if shouldApplyUpdate {
		storePeriod := utils.CalcSyncPeriod(store.FinalizedHeader.Slot)

		// Sync committee update logic
		if store.NextSyncCommitee == nil {
			store.NextSyncCommitee = update.NextSyncCommittee
		} else if updateFinalizedPeriod == storePeriod+1 {
			log.Println("sync committee updated")
			store.CurrentSyncCommitee = *store.NextSyncCommitee
			store.NextSyncCommitee = update.NextSyncCommittee
			store.PreviousMaxActiveParticipants = store.CurrentMaxActiveParticipants
			store.CurrentMaxActiveParticipants = 0
		}

		// Update finalized header
		if updateFinalizedSlot > store.FinalizedHeader.Slot {
			store.FinalizedHeader = update.FinalizedHeader

			if store.FinalizedHeader.Slot > store.OptimisticHeader.Slot {
				store.OptimisticHeader = store.FinalizedHeader
			}

			if store.FinalizedHeader.Slot%32 == 0 {
				checkpoint := toGethHeader(&store.FinalizedHeader).Hash()
				checkpointBytes := checkpoint.Bytes()
				return &checkpointBytes
			}
		}
	}

	return nil
}
func (in *Inner) apply_update(update *consensus_core.Update) {
	genUpdate := GenericUpdate{
		AttestedHeader:          update.AttestedHeader,
		SyncAggregate:           update.SyncAggregate,
		SignatureSlot:           update.SignatureSlot,
		NextSyncCommittee:       &update.NextSyncCommittee,
		NextSyncCommitteeBranch: &update.NextSyncCommitteeBranch,
		FinalizedHeader:         update.FinalizedHeader,
		FinalityBranch:          update.FinalityBranch,
	}
	checkpoint := in.apply_generic_update(&in.Store, &genUpdate)
	if checkpoint != nil {
		in.lastCheckpoint = checkpoint
	}
}
func (in *Inner) apply_finality_update(update *consensus_core.FinalityUpdate) {
	genUpdate := GenericUpdate{
		AttestedHeader:  update.AttestedHeader,
		SyncAggregate:   update.SyncAggregate,
		SignatureSlot:   update.SignatureSlot,
		FinalizedHeader: update.FinalizedHeader,
		FinalityBranch:  update.FinalityBranch,
	}
	checkpoint := in.apply_generic_update(&in.Store, &genUpdate)
	if checkpoint != nil {
		in.lastCheckpoint = checkpoint
	}
}
func (in *Inner) apply_optimistic_update(update *consensus_core.OptimisticUpdate) {
	genUpdate := GenericUpdate{
		AttestedHeader: update.AttestedHeader,
		SyncAggregate:  update.SyncAggregate,
		SignatureSlot:  update.SignatureSlot,
	}
	checkpoint := in.apply_generic_update(&in.Store, &genUpdate)
	if checkpoint != nil {
		in.lastCheckpoint = checkpoint
	}
}
func (in *Inner) Log_finality_update(update *consensus_core.FinalityUpdate) {
	participation := float32(getBits(update.SyncAggregate.SyncCommitteeBits[:])) / 512.0 * 100.0
	decimals := 2
	if participation == 100.0 {
		decimals = 1
	}

	age := in.Age(in.Store.FinalizedHeader.Slot)
	days := age.Hours() / 24
	hours := int(age.Hours()) % 24
	minutes := int(age.Minutes()) % 60
	seconds := int(age.Seconds()) % 60

	log.Printf(
		"finalized slot             slot=%d  confidence=%.*f%%  age=%02d:%02d:%02d:%02d",
		in.Store.FinalizedHeader.Slot, decimals, participation, int(days), hours, minutes, seconds,
	)
}
func (in *Inner) Log_optimistic_update(update *consensus_core.OptimisticUpdate) {
	participation := float32(getBits(update.SyncAggregate.SyncCommitteeBits[:])) / 512.0 * 100.0
	decimals := 2
	if participation == 100.0 {
		decimals = 1
	}

	age := in.Age(in.Store.OptimisticHeader.Slot)
	days := age.Hours() / 24
	hours := int(age.Hours()) % 24
	minutes := int(age.Minutes()) % 60
	seconds := int(age.Seconds()) % 60

	log.Printf(
		"updated head               slot=%d  confidence=%.*f%%  age=%02d:%02d:%02d:%02d",
		in.Store.OptimisticHeader.Slot, decimals, participation, int(days), hours, minutes, seconds,
	)
}
func (in *Inner) has_finality_update(update *GenericUpdate) bool {
	return update.FinalizedHeader != (consensus_core.Header{}) && update.FinalityBranch != nil
}
func (in *Inner) has_sync_update(update *GenericUpdate) bool {
	return update.NextSyncCommittee != nil && update.NextSyncCommitteeBranch != nil
}
func (in *Inner) safety_threshold() uint64 {
	return max(in.Store.CurrentMaxActiveParticipants, in.Store.PreviousMaxActiveParticipants) / 2
}

// verifySyncCommitteeSignature verifies the sync committee signature.
func verifySyncCommitteeSignature(
	pks []consensus_core.BLSPubKey, // Public keys slice
	attestedHeader *consensus_core.Header, // Attested header
	signature *consensus_core.SyncAggregate, // Signature bytes
	forkDataRoot consensus_core.Bytes32, // Fork data root
) bool {

	if signature == nil {
		fmt.Println("no signature")
		return false
	}

	// Create a slice to hold the collected public keys.
	collectedPks := make([]*bls.Pubkey, 0, len(pks))
	for i := range pks {
		var pksinBytes [48]byte = [48]byte(pks[i])
		dkey := new(bls.Pubkey)
		err := dkey.Deserialize(&pksinBytes)
		if err != nil {
			fmt.Println("error deserializing public key:", err)
			return false
		}

		// Include the public key only if the corresponding SyncCommitteeBits bit is set.
		if signature.SyncCommitteeBits[i/8]&(byte(1)<<(i%8)) != 0 {
			collectedPks = append(collectedPks, dkey)
		}
	}

	// Compute signingRoot.
	signingRoot := ComputeCommitteeSignRoot(toGethHeader(attestedHeader), forkDataRoot)

	var sig bls.Signature
	signatureForUnmarshalling := [96]byte(signature.SyncCommitteeSignature)

	// Deserialize the signature.
	if err := sig.Deserialize(&signatureForUnmarshalling); err != nil {
		fmt.Println("error deserializing signature:", err)
		return false
	}

	// Check if we have collected any public keys before proceeding.
	if len(collectedPks) == 0 {
		fmt.Println("no valid public keys collected")
		return false
	}

	return utils.FastAggregateVerify(collectedPks, signingRoot[:], &sig)

}

func ComputeCommitteeSignRoot(header *beacon.Header, fork consensus_core.Bytes32) consensus_core.Bytes32 {
	// Domain type for the sync committee
	domainType := [4]byte{7, 0, 0, 0}

	// Compute the domain
	domain := utils.ComputeDomain(domainType, fork)

	// Compute and return the signing root
	return utils.ComputeSigningRoot(header, domain)
}
func (in *Inner) Age(slot uint64) time.Duration {
	expectedTime := slot*12 + in.Config.Chain.GenesisTime
	now := time.Now().Unix()
	return time.Duration(uint64(now)-expectedTime) * time.Second
}
func (in *Inner) expected_current_slot() uint64 {
	const SLOT_DURATION = 12

	now := time.Now().Unix()

	sinceGenesis := now - int64(in.Config.Chain.GenesisTime)
	return uint64(sinceGenesis) / SLOT_DURATION
}
func (in *Inner) is_valid_checkpoint(blockHashSlot uint64) bool {
	const SLOT_DURATION = 12
	currentSlot := in.expected_current_slot()
	currentSlotTimestamp := int64(in.Config.Chain.GenesisTime) + int64(currentSlot)*SLOT_DURATION
	blockhashSlotTimestamp := int64(in.Config.Chain.GenesisTime) + int64(blockHashSlot)*SLOT_DURATION

	slotAge := currentSlotTimestamp - blockhashSlotTimestamp
	return uint64(slotAge) < in.Config.MaxCheckpointAge
}

func isFinalityProofValid(attestedHeader *consensus_core.Header, finalizedHeader *consensus_core.Header, finalityBranch []consensus_core.Bytes32) bool {

	return utils.IsProofValid(attestedHeader, toGethHeader(finalizedHeader).Hash(), finalityBranch, 6, 105)
}

func isCurrentCommitteeProofValid(attestedHeader *consensus_core.Header, currentCommittee *consensus_core.SyncCommittee, currentCommitteeBranch []consensus_core.Bytes32) bool {

	return utils.IsProofValid(attestedHeader, toGethSyncCommittee(currentCommittee).Root(), currentCommitteeBranch, 5, 54)
}

func isNextCommitteeProofValid(attestedHeader *consensus_core.Header, nextCommittee *consensus_core.SyncCommittee, nextCommitteeBranch []consensus_core.Bytes32) bool {

	return utils.IsProofValid(attestedHeader, toGethSyncCommittee(nextCommittee).Root(), nextCommitteeBranch, 5, 55)
}

func PayloadToBlock(value *consensus_core.ExecutionPayload) (*common.Block, error) {
	emptyNonce := "0x0000000000000000"
	emptyUncleHash := geth.HexToHash("1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347")

	// Allocate txs on the heap by using make with a predefined capacity
	txs := make([]common.Transaction, 0, len(value.Transactions))

	// Process each transaction
	for i := range value.Transactions {
		tx, err := processTransaction(&value.Transactions[i], value.BlockHash, &value.BlockNumber, uint64(i))
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	// Construct and return the block
	return &common.Block{
		Number:           value.BlockNumber,
		BaseFeePerGas:    *uint256.NewInt(value.BaseFeePerGas),
		Difficulty:       *uint256.NewInt(0),
		ExtraData:        value.ExtraData[:],
		GasLimit:         value.GasLimit,
		GasUsed:          value.GasUsed,
		Hash:             value.BlockHash,
		LogsBloom:        value.LogsBloom[:],
		ParentHash:       value.ParentHash,
		ReceiptsRoot:     value.ReceiptsRoot,
		StateRoot:        value.StateRoot,
		Timestamp:        value.Timestamp,
		TotalDifficulty:  uint64(0),
		Transactions:     common.Transactions{Full: txs},
		MixHash:          value.PrevRandao,
		Nonce:            emptyNonce,
		Sha3Uncles:       emptyUncleHash,
		Size:             0,
		TransactionsRoot: [32]byte{},
		Uncles:           [][32]byte{},
		BlobGasUsed:      value.BlobGasUsed,
		ExcessBlobGas:    value.ExcessBlobGas,
	}, nil
}

func processTransaction(txBytes *[]byte, blockHash consensus_core.Bytes32, blockNumber *uint64, index uint64) (common.Transaction, error) {
	// Decode the transaction envelope (RLP-encoded)

	txEnvelope, err := DecodeTxEnvelope(txBytes)
	if err != nil {
		return common.Transaction{}, fmt.Errorf("failed to decode transaction: %v", err)
	}
	// Updated due to update in Transaction struct
	tx := common.Transaction{
		Hash:             txEnvelope.Hash(),
		Nonce:            hexutil.Uint64(txEnvelope.Nonce()),
		BlockHash:        geth.BytesToHash(blockHash[:]),
		BlockNumber:      hexutil.Uint64(*blockNumber),
		TransactionIndex: hexutil.Uint64(index),
		To:               txEnvelope.To(),
		Value:            hexutil.Big(*txEnvelope.Value()),
		GasPrice:         hexutil.Big(*txEnvelope.GasPrice()),
		Gas:              hexutil.Uint64(txEnvelope.Gas()),
		Input:            txEnvelope.Data(),
		ChainID:          hexutil.Big(*txEnvelope.ChainId()),
		TransactionType:  hexutil.Uint(uint(txEnvelope.Type())),
	}

	// Handle signature and transaction type logic
	signer := types.LatestSignerForChainID(txEnvelope.ChainId())
	from, err := types.Sender(signer, txEnvelope)
	if err != nil {
		return common.Transaction{}, fmt.Errorf("failed to recover sender: %v", err)
	}
	tx.From = &from

	// Extract signature components
	r, s, v := txEnvelope.RawSignatureValues()
	tx.Signature = &common.Signature{
		R:       r.String(),
		S:       s.String(),
		V:       v.Uint64(),
		YParity: common.Parity{Value: v.Uint64() == 1},
	}

	switch txEnvelope.Type() {
	case types.AccessListTxType:
		tx.AccessList = txEnvelope.AccessList()
	case types.DynamicFeeTxType:
		// Update due to update in Transaction struct
		tx.MaxFeePerGas = hexutil.Big(*new(big.Int).Set(txEnvelope.GasFeeCap()))
		tx.MaxPriorityFeePerGas = hexutil.Big(*new(big.Int).Set(txEnvelope.GasTipCap()))
	case types.BlobTxType:
		tx.MaxFeePerGas = hexutil.Big(*new(big.Int).Set(txEnvelope.GasFeeCap()))
		tx.MaxPriorityFeePerGas = hexutil.Big(*new(big.Int).Set(txEnvelope.GasTipCap()))
		tx.MaxFeePerBlobGas = hexutil.Big(*new(big.Int).Set(txEnvelope.BlobGasFeeCap()))
		tx.BlobVersionedHashes = txEnvelope.BlobHashes()
	case types.LegacyTxType:
		// No additional fields to set
	default:
		fmt.Print("Unhandled transaction type")
	}

	return tx, nil
}

// getBits counts the number of bits set to 1 in a [64]byte array
func getBits(bitfield []byte) uint64 {
	var count uint64
	for _, b := range bitfield {
		count += uint64(popCount(b))
	}
	return count
}

// popCount counts the number of bits set to 1 in a byte
func popCount(b byte) int {
	count := 0
	for b != 0 {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

// DecodeTxEnvelope takes the transaction bytes and decodes them into a transaction envelope (Ethereum transaction)
func DecodeTxEnvelope(txBytes *[]byte) (*types.Transaction, error) {
	// Create an empty transaction object
	var tx types.Transaction

	var txBytesForUnmarshal []byte = *txBytes
	// Unmarshal the RLP-encoded transaction bytes into the transaction object
	err := tx.UnmarshalBinary(txBytesForUnmarshal)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %v", err)
	}

	return &tx, nil
}

func SomeGasPrice(gasFeeCap, gasTipCap *big.Int, baseFeePerGas uint64) *big.Int {
	// Create a new big.Int for the base fee and set its value
	baseFee := new(big.Int).SetUint64(baseFeePerGas)

	// Calculate the maximum gas price based on the provided parameters
	maxGasPrice := new(big.Int).Set(gasFeeCap)

	// Calculate the alternative gas price
	alternativeGasPrice := new(big.Int).Set(gasTipCap)
	alternativeGasPrice.Add(alternativeGasPrice, baseFee)

	// Return the maximum of the two
	if maxGasPrice.Cmp(alternativeGasPrice) < 0 {
		return alternativeGasPrice
	}
	return maxGasPrice
}

func toGethHeader(header *consensus_core.Header) *beacon.Header {
	return &beacon.Header{
		Slot:          header.Slot,
		ProposerIndex: header.ProposerIndex,
		ParentRoot:    [32]byte(header.ParentRoot),
		StateRoot:     [32]byte(header.StateRoot),
		BodyRoot:      [32]byte(header.BodyRoot),
	}
}

type jsonSyncCommittee struct {
	Pubkeys   []hexutil.Bytes
	Aggregate hexutil.Bytes
}

func toGethSyncCommittee(committee *consensus_core.SyncCommittee) *beacon.SerializedSyncCommittee {
	jsoncommittee := &jsonSyncCommittee{
		Aggregate: committee.AggregatePubkey[:],
	}

	for _, pubkey := range committee.Pubkeys {
		jsoncommittee.Pubkeys = append(jsoncommittee.Pubkeys, pubkey[:])
	}

	var s beacon.SerializedSyncCommittee

	for i, key := range jsoncommittee.Pubkeys {
		copy(s[i*48:], key[:])
	}
	copy(s[512*48:], jsoncommittee.Aggregate[:])
	return &s
}

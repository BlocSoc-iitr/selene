package consensus

///NOTE: only these imports are required others are in package already
// uses rpc
// uses config for networks
// uses common for datatypes
type ConsensusClient struct{}
type Inner struct{}
type LightClientStore struct{}

func(con ConsensusClient) new(){}
func(con ConsensusClient) shutdown(){}
func(con ConsensusClient) expected_current_slot(){}

func sync_fallback(){}
func sync_all_fallback(){}


func(in Inner) new(){}
func(in Inner) get_rpc(){}
func(in Inner) check_execution_payload(){}
func(in Inner) get_payloads(){}
func(in Inner) advance(){}
func(in Inner) send_blocks(){}
func(in Inner) duration_until_next_update(){}
func(in Inner) bootstrap(){}
func(in Inner) verify_generic_update(){}
func(in Inner) verify_update(){}
func(in Inner) verify_finality_update(){}
func(in Inner) verify_optimistic_update(){}
func(in Inner) apply_generic_update(){}
func(in Inner) apply_update(){}
func(in Inner) apply_finality_update(){}
func(in Inner) apply_optimistic_update(){}
func(in Inner) log_finality_update(){}
func(in Inner) log_optimistic_update(){}
func(in Inner) has_finality_update(){}
func(in Inner) has_sync_update(){}
func(in Inner) safety_threshold(){}
func(in Inner) verify_sync_committee_signature(){}
func(in Inner) compute_committee_sign_root(){}
func(in Inner) age(){}
func(in Inner) expected_current_slot(){}
func(in Inner) is_valid_checkpoint(){}



func get_participating_keys(){}
func is_finality_proof_valid(){}
func is_current_committee_proof_valid(){}
func is_next_committee_proof_valid(){}
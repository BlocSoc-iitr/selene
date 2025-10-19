package evm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestSStoreResult(t *testing.T){
	result :=SStoreResult{
		OriginalValue: U256(big.NewInt(10)),
		PresentValue: U256(big.NewInt(20)),
		NewValue: U256(big.NewInt(30)),
		IsCold: true,
		}
	
		assert.Equal(t, result.OriginalValue, big.NewInt(10), "Original Value should be equal to 10")
		assert.Equal(t, result.PresentValue, big.NewInt(20), "Present Value should be equal to 20")
		assert.Equal(t, result.NewValue, big.NewInt(30), "New Value should be equal to 30")
		assert.True(t, result.IsCold, "IsCold should be true")

}
func TestLoadAccountResult(t *testing.T){
	result:=LoadAccountResult{
		IsCold: true,
		IsEmpty: false,

	}
	assert.True(t,result.IsCold,"IsCold should be true")
	assert.False(t,result.IsEmpty,"IsEmpty should be true")
}
func TestSelfDestrcutResult(t *testing.T){
	result:=SelfDestructResult{
		HadValue: true,
		TargetExists: false,
		IsCold: false,
		PreviouslyDestroyed: true,
	}
	assert.True(t,result.HadValue,"HardValue should be true")
	assert.False(t,result.TargetExists,"TargetExists should be false")
	assert.False(t,result.IsCold,"IsCold should be false")
	assert.True(t,result.PreviouslyDestroyed,"Previously Deployed should be true")
}

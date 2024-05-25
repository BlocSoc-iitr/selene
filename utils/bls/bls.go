package bls

import (
	"math/big"
	"github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fp"
)

type G1Point struct {
	*bls12381.G1Affine
}

type G2Point struct {
	*bls12381.G2Affine
}

type Signature struct {
	*G1Point `json:"g1_point"`
}

// TODO: this function uses G1 for sigs and G2 for pubkeys, however
// ethereum consensus uses it vice-versa. Research more about this.
func (s *Signature) Verify(pubkey *G2Point, message [32]byte) (bool, error) {
	g2Gen := GetG2Generator()

	msgPoint := MapToCurve(message)

	var negSig bls12381.G1Affine
	negSig.Neg((*bls12381.G1Affine)(s.G1Affine))

	P := [2]bls12381.G1Affine{*msgPoint, negSig}
	Q := [2]bls12381.G2Affine{*pubkey.G2Affine, *g2Gen}

	ok, err := bls12381.PairingCheck(P[:], Q[:])
	if err != nil {
		return false, err
	}
	return ok, nil
}

// TODO: change function to comply with bls12381
// bls function: y^2 = x^3 + 4
func MapToCurve(digest [32]byte) *bls12381.G1Affine {

	one := new(big.Int).SetUint64(1)
	three := new(big.Int).SetUint64(3)
	x := new(big.Int)
	x.SetBytes(digest[:])
	for {
		// y = x^3 + 3
		xP3 := new(big.Int).Exp(x, big.NewInt(3), fp.Modulus())
		y := new(big.Int).Add(xP3, three)
		y.Mod(y, fp.Modulus())

		if y.ModSqrt(y, fp.Modulus()) == nil {
			x.Add(x, one).Mod(x, fp.Modulus())
		} else {
			var fpX, fpY fp.Element
			fpX.SetBigInt(x)
			fpY.SetBigInt(y)
			return &bls12381.G1Affine{
				X: fpX,
				Y: fpY,
			}
		}
	}
}

// TODO: change the values from bn254 to bls12381
func GetG2Generator() *bls12381.G2Affine {
	g2Gen := new(bls12381.G2Affine)
	g2Gen.X.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781",
		"11559732032986387107991004021392285783925812861821192530917403151452391805634")
	g2Gen.Y.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930",
		"4082367875863433681332203403145435568316851327593401208105741076214120093531")
	return g2Gen
}
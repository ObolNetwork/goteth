package spec

import "github.com/attestantio/go-eth2-client/spec/phase0"

// https://github.com/attestantio/vouch/blob/0c75ee8315dc4e5df85eb2aa09b4acc2b4436661/strategies/beaconblockproposal/best/score.go#L426
// intersection returns a list of items common between the two sets.
func SlashingIntersection(set1 []uint64, set2 []uint64) []phase0.ValidatorIndex {

	res := make([]phase0.ValidatorIndex, 0)
	for _, item1 := range set1 {
		for _, item2 := range set2 {
			if item1 == item2 {
				res = append(res, phase0.ValidatorIndex(item1))
			}
		}
	}

	return res

}

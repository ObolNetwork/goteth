package analyzer

import "github.com/attestantio/go-eth2-client/spec/phase0"

// setEpochBoundaryStateRoot stores the state root from a Head SSE event
// for the given epoch-boundary slot.
func (s *ChainAnalyzer) setEpochBoundaryStateRoot(slot phase0.Slot, root phase0.Root) {
	s.epochBoundaryStateRoots.Store(slot, root)
}

// takeEpochBoundaryStateRoot retrieves and deletes the cached state root
// for the given slot. Returns the root and true if found, or a zero root
// and false if not cached.
func (s *ChainAnalyzer) takeEpochBoundaryStateRoot(slot phase0.Slot) (phase0.Root, bool) {
	val, ok := s.epochBoundaryStateRoots.LoadAndDelete(slot)
	if !ok {
		return phase0.Root{}, false
	}
	return val.(phase0.Root), true
}

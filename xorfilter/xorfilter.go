package xorfilter

import (
	"bufio"
	"container/list"
	"encoding/binary"
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/cespare/xxhash/v2"
)

// --------------------------------------------------------------------------------------------- //
// Cell represents a single bit set in the XOR filter
// --------------------------------------------------------------------------------------------- //
type Cell struct {
	Bits uint64 // (L = 64 bits, error rate ~ 1/2^64)
}

// --------------------------------------------------------------------------------------------- //
// XORfilter is a probabilistic data structure for set membership queries
// --------------------------------------------------------------------------------------------- //
type XORfilter struct {
	Elements  []string // Keys used to build the filter (cleared after building)
	Size      int      // Total size of the filter, calculated as floor(1.23 * len(Elements)) + 32
	BlockSize int      // Size of each block in the filter
	Seed      uint64   // Seed for hash functions to ensure reproducibility
	BitSet    []Cell   // Array of cells storing the fingerprints
}

// --------------------------------------------------------------------------------------------- //
// Mapping holds a key and its corresponding hash positions
// and fingerprint during filter construction
// --------------------------------------------------------------------------------------------- //
type Mapping struct {
	Key         []byte // The original key
	H0, H1, H2  uint64 // Positions in the three blocks
	Fingerprint uint64 // Computed fingerprint of the key
}

// --------------------------------------------------------------------------------------------- //
// ProcessKeys reads keys from standard input or a file specified in command-line arguments.
// It returns a pointer to a slice of keys or an error if reading fails.
// --------------------------------------------------------------------------------------------- //
func ProcessKeys() (*[]string, error) {
	var stream *os.File
	var err error
	var keys []string

	if len(os.Args) < 2 {
		stream = os.Stdin
	} else {
		source := os.Args[1]
		stream, err = os.Open(source)
		if err != nil {
			return nil, err
		}
		defer stream.Close()
	}

	input := bufio.NewScanner(stream)
	for input.Scan() {
		keys = append(keys, input.Text())
	}
	if err := input.Err(); err != nil {
		return nil, err
	}

	return &keys, nil
}

// --------------------------------------------------------------------------------------------- //
// InitializeXORfilter creates a new XOR filter for the given set of keys.
// It returns a pointer to the initialized filter.
// --------------------------------------------------------------------------------------------- //
func InitializeXORfilter(keys []string) *XORfilter {
	size := int(float64(len(keys))*1.23) + 32
	if size < 3 {
		size = 3 // Minimum size
	}

	seed := uint64(rand.NewSource(time.Now().UnixNano()).Int63())
	bitset := make([]Cell, size)

	return &XORfilter{
		Elements:  keys,
		Size:      size,
		BlockSize: size / 3,
		Seed:      seed,
		BitSet:    bitset,
	}
}

// --------------------------------------------------------------------------------------------- //
// computeHash computes a hash value for the given key, seed, and modifier.
// --------------------------------------------------------------------------------------------- //
func computeHash(key []byte, seedBytes []byte, modifier byte) uint64 {
	// h := fnv.New64a()
	h := xxhash.New()

	h.Write(key)
	h.Write(seedBytes)
	h.Write([]byte{modifier})

	return h.Sum64()
}

// --------------------------------------------------------------------------------------------- //
// getHashes computes three hash values for a key,
// each mapped to a position in one of the three blocks.
// It returns an array of three positions.
// --------------------------------------------------------------------------------------------- //
func getHashes(filter *XORfilter, key []byte, seed uint64) [3]uint64 {
	blockSize := filter.BlockSize
	block0End := blockSize
	block1End := 2 * blockSize

	seedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(seedBytes, seed)

	hash0 := computeHash(key, seedBytes, 0)
	hash1 := computeHash(key, seedBytes, 1)
	hash2 := computeHash(key, seedBytes, 2)

	return [3]uint64{
		hash0 % uint64(blockSize),
		uint64(block0End) + (hash1 % uint64(blockSize)),
		uint64(block1End) + (hash2 % uint64(filter.Size-block1End)),
	}
}

// --------------------------------------------------------------------------------------------- //
// getFingerprintHash computes the fingerprint hash for a key.
// It returns the fingerprint as a uint64.
// --------------------------------------------------------------------------------------------- //
func getFingerprintHash(key []byte, seed uint64) uint64 {
	seedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(seedBytes, seed)
	return computeHash(key, seedBytes, 3)
}

// --------------------------------------------------------------------------------------------- //
// tryBuild attempts to build the XOR filter with the current seed.
// It returns true if successful, false if not all keys were mapped.
// --------------------------------------------------------------------------------------------- //
func (filter *XORfilter) tryBuild() bool {
	if len(filter.Elements) == 0 {
		return false
	}

	// Prepare mappings for each key
	mappings := make([]Mapping, len(filter.Elements))
	for i, key := range filter.Elements {
		keyBytes := []byte(key)
		h := getHashes(filter, keyBytes, filter.Seed)
		mappings[i] = Mapping{
			Key:         keyBytes,
			H0:          h[0],
			H1:          h[1],
			H2:          h[2],
			Fingerprint: getFingerprintHash(keyBytes, filter.Seed),
		}
	}

	// Count how many keys map to each position
	count := make([]int, filter.Size)
	for _, m := range mappings {
		count[m.H0]++
		count[m.H1]++
		count[m.H2]++
	}

	// Initialize queue with singleton positions
	queue := list.New()
	for i, c := range count {
		if c == 1 {
			queue.PushBack(i)
		}
	}

	// Stack entries record order of removal
	type peelEntry struct{ idx, pos int }
	stack := make([]peelEntry, 0, len(mappings))
	used := make([]bool, len(mappings))

	// Peeling process
	for queue.Len() > 0 {
		e := queue.Front()
		pos := e.Value.(int)
		queue.Remove(e)

		// Find mapping covering this pos
		for i := range mappings {
			if !used[i] && (mappings[i].H0 == uint64(pos) || mappings[i].H1 == uint64(pos) || mappings[i].H2 == uint64(pos)) {
				used[i] = true
				stack = append(stack, peelEntry{idx: i, pos: pos})

				// Decrement neighbor counts
				for _, nh := range []uint64{mappings[i].H0, mappings[i].H1, mappings[i].H2} {
					count[nh]--
					if count[nh] == 1 {
						queue.PushBack(int(nh))
					}
				}
				break
			}
		}
	}

	// Check if all keys peeled
	for _, u := range used {
		if !u {
			return false
		}
	}

	// Assign fingerprints in reverse peel order
	for i := len(stack) - 1; i >= 0; i-- {
		entry := stack[i]
		m := mappings[entry.idx]
		mask := m.Fingerprint

		mask ^= filter.BitSet[m.H0].Bits
		mask ^= filter.BitSet[m.H1].Bits
		mask ^= filter.BitSet[m.H2].Bits

		filter.BitSet[entry.pos].Bits = mask
	}

	return true
}

// --------------------------------------------------------------------------------------------- //
// Build constructs the XOR filter by assigning fingerprints to BitSet.
// It retries with different seeds until successful or max attempts are reached.
// It returns an error if the construction ultimately fails.
// --------------------------------------------------------------------------------------------- //
func (filter *XORfilter) Build() error {
	const maxAttempts = 1000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for attempt := 0; attempt < maxAttempts; attempt++ {
		filter.Seed = uint64(r.Int63())

		if filter.tryBuild() {
			filter.Elements = nil // Clear keys to save memory
			return nil
		}
	}

	return errors.New("failed to build XOR filter after multiple attempts")
}

// --------------------------------------------------------------------------------------------- //
// Contains checks if a key is likely present in the filter.
// It returns true if the key might be in the set (with possible false positives).
// --------------------------------------------------------------------------------------------- //
func (filter *XORfilter) Contains(key []byte) bool {
	hashes := getHashes(filter, key, filter.Seed)
	fingerprint := getFingerprintHash(key, filter.Seed)

	return filter.BitSet[hashes[0]].Bits^
		filter.BitSet[hashes[1]].Bits^
		filter.BitSet[hashes[2]].Bits == fingerprint
}

// --------------------------------------------------------------------------------------------- //

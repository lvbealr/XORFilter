package main

import (
	"XORFilter/xorfilter"
	"fmt"
	"os"
)

func main() {
	// Read keys from stdin or from file
	keys, err := xorfilter.ProcessKeys()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading keys: %v\n", err)
		os.Exit(1)
	}

	// Create XOR filter
	filter := xorfilter.InitializeXORfilter(*keys)
	fmt.Printf("Initialized XOR filter with %d keys, size %d\n", len(*keys), filter.Size)

	// Build filter
	if err := filter.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "Error building filter: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Filter built successfully")

	// Test existing keys
	for _, key := range *keys {
		if !filter.Contains([]byte(key)) {
			fmt.Fprintf(os.Stderr, "Error: key %q should be in the filter\n", key)
		} else {
			fmt.Printf("Key %q is likely in the filter\n", key)
		}
	}

	// Test nonexistent key
	testKey := []byte("nonexistent_key")
	if filter.Contains(testKey) {
		fmt.Printf("Key %q is likely in the filter (false positive)\n", testKey)
	} else {
		fmt.Printf("Key %q is not in the filter\n", testKey)
	}
}

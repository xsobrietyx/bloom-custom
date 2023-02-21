package bloom_custom

import (
	"github.com/shivakar/metrohash"
	"github.com/spaolacci/murmur3"
	"log"
	"math"
)

type Bloom interface {
	New(size float64, probability float64) Bloom
	Set(value []byte)
	Verify(value []byte) bool
}
type Filter struct {
	internal  []byte
	functions []func(v []byte) uint64
}

func murmur3HashFunction(v []byte) uint64 {
	hash := murmur3.New64()
	_, err := hash.Write(v)
	if err != nil {
		panic(err)
	}
	return hash.Sum64()
}

func metroHashFunction(v []byte) uint64 {
	hash := metrohash.NewMetroHash64()
	_, err := hash.Write(v)
	if err != nil {
		panic(err)
	}
	return hash.Sum64()
}

func initHashes(functionsAmount uint) []func(v []byte) uint64 {
	//TODO: for larger data sets additional hash functions should be added
	hashFunctions := make([]func(v []byte) uint64, 10)
	hashFunctions[0] = murmur3HashFunction
	hashFunctions[1] = metroHashFunction
	hashFunctions = hashFunctions[0:functionsAmount]
	return hashFunctions
}

func (f *Filter) New(size float64, prob float64) Bloom {
	logs := log.Default()
	logs.SetPrefix("[bloom-custom]")

	bytesBufferSize := math.Round(-(size * math.Log(prob)) / (math.Ln2 * math.Ln2))
	hashFunctionsAmount := math.Round((bytesBufferSize / size) * math.Ln2)

	log.Printf("Buffer size: %v, size: %v, # of hash functions used: %v\n", bytesBufferSize, size, hashFunctionsAmount)
	return &Filter{internal: make([]byte, uint64(bytesBufferSize),
		uint64(bytesBufferSize)),
		functions: initHashes(uint(hashFunctionsAmount))}
}

func (f *Filter) Set(value []byte) {
	for _, v := range f.functions {
		foo := v
		index := foo(value) % uint64(len(f.internal))
		if index < 0 {
			index = -index
		}
		f.internal[index] = 1
	}
}

func (f *Filter) Verify(value []byte) bool {
	res := false
	for _, v := range f.functions {
		foo := v
		index := foo(value) % uint64(len(f.internal))
		if index < 0 {
			index = -index
		}
		if f.internal[index] == 1 {
			res = true
		} else {
			res = false
		}
	}
	return res
}

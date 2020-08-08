package kerlong

import (
	"github.com/willf/bitset"
)

type SimpleHash struct {
	cap  uint
	seed uint
}

type BloomFilter struct {
	Size  uint
	Set   *bitset.BitSet
	Funcs []SimpleHash
}

func (s SimpleHash) Hash(value string) uint {
	var result uint = 0
	for i := 0; i < len(value); i++ {
		result = result*s.seed + uint(value[i])
	}
	return (s.cap - 1) & result
}

//Contains 是否可能存在
func (bf BloomFilter) Contains(value string) bool {
	if value == "" {
		return false
	}
	ret := true
	for _, f := range bf.Funcs {
		ret = ret && bf.Set.Test(f.Hash(value))
	}
	return ret
}

//NewBloomFilter 新建一个bloomfilter
func NewBloomFilter(size uint, seeds []uint) *BloomFilter {
	bf := new(BloomFilter)
	if size < uint(len(seeds)) {
		bf.Size = 2 << 24

	} else {
		bf.Size = size
	}
	bf.Funcs = make([]SimpleHash, len(seeds))
	for i := 0; i < len(bf.Funcs); i++ {
		bf.Funcs[i] = SimpleHash{bf.Size, seeds[i]}
	}
	bf.Set = bitset.New(bf.Size)
	return bf
}

//Add 添加一个数据
func (bf BloomFilter) Add(value string) {
	for _, f := range bf.Funcs {
		bf.Set.Set(f.Hash(value))
	}
}

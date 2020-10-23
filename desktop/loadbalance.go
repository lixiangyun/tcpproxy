package main

import (
	"math/rand"
	"net"
	"strings"
	"sync/atomic"
)

type LoadBalance interface {
	Next(addr string) int
}

type RandomLB struct {
	max  int
}

func (r *RandomLB)Next(addr string) int {
	return rand.Int() % r.max
}

func NewRandom(items []BackendConfig) LoadBalance {
	return &RandomLB{max: len(items)}
}

type RoundRobinLB struct {
	idx uint32
	max uint32
}

func (r *RoundRobinLB)Next(addr string) int {
	idx := r.idx
	atomic.AddUint32(&r.idx, 1)
	return int(idx % r.max)
}

func NewRoundRobin(items []BackendConfig) LoadBalance {
	return &RoundRobinLB{max: uint32(len(items))}
}

type AddressHashLB struct {
	max int
}

func (r *AddressHashLB)Next(addr string) int {
	idx := strings.Index(addr, ":")
	if idx != -1 {
		addr = addr[:idx]
	}
	ip := net.ParseIP(addr)
	var sum int
	for _, v := range ip {
		sum += int(v)
	}
	return sum % r.max
}

func NewAddressHash(items []BackendConfig) LoadBalance {
	return &AddressHashLB{max: len(items)}
}

type WeightRoundRobinLB struct {
	List  []int
	Idx     int
	Gcd     int
	MaxW    int
	CurW    int
}

func (r *WeightRoundRobinLB)Next(addr string) int {
	for {
		r.Idx = (r.Idx + 1) % len(r.List)
		if r.Idx == 0 {
			r.CurW = r.CurW - r.Gcd
			if r.CurW <= 0 {
				r.CurW = r.MaxW
			}
		}
		if r.List[r.Idx] >= r.CurW {
			return r.Idx
		}
	}
}

func NewWeightRoundRobin(items []BackendConfig) LoadBalance {
	s := &WeightRoundRobinLB{List: make([]int, len(items))}
	for idx, v := range items {
		s.List[idx] = v.Weight
		if s.MaxW < v.Weight {
			s.MaxW = v.Weight
		}
		if idx == 0 {
			s.Gcd = v.Weight
		} else {
			s.Gcd = gcd(s.Gcd, v.Weight)
		}
	}
	s.CurW = s.MaxW
	return s
}

/* 迭代法（递推法）：欧几里得算法，计算最大公约数 */
func gcd(m, n int) int {
	for {
		if m == 0 {
			return n
		}
		c := n % m
		n = m
		m = c
	}
}

func NewLoadBalance(mode string, items []BackendConfig) LoadBalance {
	switch mode {
	case "Random":
		return NewRandom(items)
	case "RoundRobin":
		return NewRoundRobin(items)
	case "WeightRoundRobin":
		return NewWeightRoundRobin(items)
	case "AddressHash":
		return NewAddressHash(items)
	case "MainStandby":
		return nil
	}
	return nil
}
/*  Copyright (c) 2014, Brian Hummer (brian@boggo.net)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of the boggo.net nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL BRIAN HUMMER BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package random

import (
	"math"
	"math/rand"
	"time"
)

var (
	rng *Random
)

func init() {
	rng = New(time.Now().UnixNano())
}

func Reseed(seed int64) {
	if rng != nil && rng.running {
		rng.Close()
	}
	rng = New(seed)
}

func Next() float64 {
	return rng.Next()
}

func Gaussian() float64 {
	return rng.Gaussian()
}

func Int(n int) int {
	return rng.Int(n)
}

// Random number generator with buffered channel backends
type Random struct {
	rng     *rand.Rand
	floats  chan float64
	gauss   chan float64
	running bool
	iset    bool
	gset    float64
}

// Constucts a new Random
func New(seed int64) *Random {
	r := &Random{rng: rand.New(rand.NewSource(seed)),
		floats:  make(chan float64, 20),
		gauss:   make(chan float64, 20),
		running: true, iset: false}
	go r.processFloats()
	go r.processGauss()
	return r
}

func (r *Random) processFloats() {
	defer func() { recover() }()
	for r.running {
		r.floats <- r.rng.Float64()
	}
}

func (r *Random) processGauss() {
	defer func() { recover() }()
	for r.running {
		r.gauss <- r.nextGauss()
	}
}

func (r *Random) nextGauss() float64 {
	var fac, rsq, v1, v2 float64
	if r.iset == false {
		rsq = 0
		for rsq >= 1.0 || rsq == 0.0 {
			v1 = 2.0*r.Next() - 1.0
			v2 = 2.0*r.Next() - 1.0
			rsq = v1*v1 + v2*v2
		}
		fac = math.Sqrt(-2.0 * math.Log(rsq) / rsq)
		r.gset = v1 * fac
		r.iset = true
		return v2 * fac
	} else {
		r.iset = false
		return r.gset
	}
}

func (r *Random) Close() {
	r.running = false
	close(r.floats)
	close(r.gauss)
}

func (r *Random) Next() float64 {
	return <-r.floats
}

func (r *Random) Int(n int) int {
	return int(r.Next() * float64(n))
}

func (r *Random) Gaussian() float64 {
	return <-r.gauss
}

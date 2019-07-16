// Copyright (c) 2019 Kevin Kragenbrink, II
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package services

import (
	"math/rand"

	"github.com/kkragenbrink/slate/settings"
	"github.com/sgade/randomorg"
)

type mode int8

const (
	modeMathRandom mode = 1
	modeRandomOrg  mode = 2
)

// The RandomService is responsible for generating random numbers
type RandomService struct {
	settings *settings.Settings
	mode     mode
	random   *randomorg.Random
}

// NewRandom instantiates and configures the random service
func NewRandom(set *settings.Settings) *RandomService {
	rand := new(RandomService)
	rand.settings = set
	rand.mode = modeMathRandom

	if rand.settings.RandomOrgAPIKey != "" {
		rand.mode = modeRandomOrg
		rand.random = randomorg.NewRandom(rand.settings.RandomOrgAPIKey)
	}

	return rand
}

// Rand generates a slice of random numbers between min and max.
func (r *RandomService) Rand(times int, min, max int64) ([]int64, error) {
	if r.mode == modeRandomOrg {
		rolls, err := r.random.GenerateIntegers(times, min, max)
		if err != nil {
			return nil, err
		}
		return rolls, nil
	}

	return mathRand(times, min, max), nil
}

func mathRand(times int, min, max int64) []int64 {
	var results []int64
	for i := 0; i < times; i++ {
		roll := rand.Int63n(max) + min
		results = append(results, roll)
	}
	return results
}

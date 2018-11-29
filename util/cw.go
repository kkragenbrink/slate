// Copyright (c) 2018 Kevin Kragenbrink, II
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

package util

import (
	"fmt"
	"io"
	"os"
)

var (
	// NBlack is normal black
	NBlack = []byte{'\033', '[', '3', '0', 'm'}
	// NRed is normal red
	NRed = []byte{'\033', '[', '3', '1', 'm'}
	// NGreen is normal green
	NGreen = []byte{'\033', '[', '3', '2', 'm'}
	// NYellow is normal yellow
	NYellow = []byte{'\033', '[', '3', '3', 'm'}
	// NBlue is normal blue
	NBlue = []byte{'\033', '[', '3', '4', 'm'}
	// NMagenta is normal magenta
	NMagenta = []byte{'\033', '[', '3', '5', 'm'}
	// NCyan is normal cyan
	NCyan = []byte{'\033', '[', '3', '6', 'm'}
	//NWhite is normal white
	NWhite = []byte{'\033', '[', '3', '7', 'm'}

	// BBlack is bright black
	BBlack = []byte{'\033', '[', '3', '0', ';', '1', 'm'}
	// BRed is bright red
	BRed = []byte{'\033', '[', '3', '1', ';', '1', 'm'}
	// BGreen is bright green
	BGreen = []byte{'\033', '[', '3', '2', ';', '1', 'm'}
	// BYellow is bright yellow
	BYellow = []byte{'\033', '[', '3', '3', ';', '1', 'm'}
	// BBlue is bright blue
	BBlue = []byte{'\033', '[', '3', '4', ';', '1', 'm'}
	// BMagenta is bright magenta
	BMagenta = []byte{'\033', '[', '3', '5', ';', '1', 'm'}
	// BCyan is brihgt cyan
	BCyan = []byte{'\033', '[', '3', '6', ';', '1', 'm'}
	// BWhite is bright white
	BWhite = []byte{'\033', '[', '3', '7', ';', '1', 'm'}

	reset = []byte{'\033', '[', '0', 'm'}
)

var isTTY bool

func init() {
	// This is sort of cheating: if stdout is a character device, we assume
	// that means it's a TTY. Unfortunately, there are many non-TTY
	// character devices, but fortunately stdout is rarely set to any of
	// them.
	//
	// We could solve this properly by pulling in a dependency on
	// code.google.com/p/go.crypto/ssh/terminal, for instance, but as a
	// heuristic for whether to print in color or in black-and-white, I'd
	// really rather not.
	fi, err := os.Stdout.Stat()
	if err == nil {
		m := os.ModeDevice | os.ModeCharDevice
		isTTY = fi.Mode()&m == m
	}
}

// CW is a colorwriter
func CW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	if isTTY && useColor {
		w.Write(color)
	}
	fmt.Fprintf(w, s, args...)
	if isTTY && useColor {
		w.Write(reset)
	}
}

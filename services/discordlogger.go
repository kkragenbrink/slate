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

package services

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"github.com/kkragenbrink/slate/util"
	"log"
	"os"
	"strings"
	"time"
)

// A Severity describes the severity of a log event.
type Severity int

const (
	// LogInfo severity is for info-level events
	LogInfo Severity = iota
	// LogWarn severity is for warn-level events
	LogWarn
	// LogError severity is for error-level events
	LogError
)

// SlateLogger implements a logger
type SlateLogger struct {
	Logger *log.Logger
}

// NewSlateLogger creates a new SlateLogger instance
func NewSlateLogger() *SlateLogger {
	logger := new(SlateLogger)
	logger.Logger = log.New(os.Stdout, "", log.LstdFlags)
	return logger
}

// LogEntry describes a log entry that is being written
type LogEntry interface {
	Write(Severity, int, time.Duration)
}

type discordLogEntry struct {
	*SlateLogger
	buf *bytes.Buffer
}

// NewDiscordLogEntry is a new LogEntry specific to discord messages.
func (dl *SlateLogger) NewDiscordLogEntry(msg *discordgo.MessageCreate, cmd string, fields []string) LogEntry {
	entry := new(discordLogEntry)
	entry.SlateLogger = dl
	buffer := new(bytes.Buffer)
	entry.buf = buffer

	util.CW(entry.buf, true, util.NCyan, "\"")
	util.CW(entry.buf, true, util.NMagenta, "%s ", cmd)
	util.CW(entry.buf, true, util.NCyan, "%s\"", strings.Join(fields, " "))
	entry.buf.WriteString(" from ")
	entry.buf.WriteString(msg.Author.String())
	entry.buf.WriteString(" - ")

	return entry
}

// Write writes out the Discord log entry to the cli.
func (entry *discordLogEntry) Write(sev Severity, bytes int, elapsed time.Duration) {
	switch {
	case sev == LogInfo:
		util.CW(entry.buf, true, util.BGreen, "INFO")
	case sev == LogWarn:
		util.CW(entry.buf, true, util.BYellow, "WARN")
	case sev == LogError:
		util.CW(entry.buf, true, util.BRed, "ERR")
	}
	util.CW(entry.buf, true, util.BBlue, " %dB", bytes)
	entry.buf.WriteString(" in ")
	if elapsed < 500*time.Millisecond {
		util.CW(entry.buf, true, util.NGreen, "%s", elapsed)
	} else if elapsed < 5*time.Second {
		util.CW(entry.buf, true, util.NYellow, "%s", elapsed)
	} else {
		util.CW(entry.buf, true, util.NRed, "%s", elapsed)
	}

	entry.Logger.Print(entry.buf.String())
}

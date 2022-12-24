// Package yapr (Yet Another Proc Reader) provides miscellaneous
// functions and types to parse /proc directory entries on Linux, with
// an emphasis on correctness.
package yapr

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Stat represents a Linux process status parsed from
// /proc/[pid]/stat. See the proc man-page from section 5 for more
// information (https://man7.org/linux/man-pages/man5/proc.5.html).
type Stat struct {
	PID                  int32
	Comm                 string
	State                rune
	PPID                 int32
	PGID                 int32
	SessionID            int32
	TTYNumber            int32
	TPGID                int32
	Flags                uint32
	MinorFaults          uint64
	ChildMinorFaults     uint64
	MajorFaults          uint64
	ChildMajorFaults     uint64
	UserTime             uint64
	SystemTime           uint64
	ChildUserTime        int64
	ChildSystemTime      int64
	Priority             int64
	Nice                 int64
	NumThreads           int64
	ITRealValue          int64
	StartTime            uint64
	VirtualMemSize       uint64
	ResidentSetSize      int64
	ResidentSetSizeLimit uint64
	StartCode            uint64
	EndCode              uint64
	StartStack           uint64
	KStackESP            uint64
	KStackEIP            uint64
	Signal               uint64
	Blocked              uint64
	SigIgnore            uint64
	SigCatch             uint64
	WChan                uint64
	NumSwap              uint64
	CumNumSwap           uint64
	ExitSignal           int32
	Processor            int32
	RTPrio               uint32
	Policy               uint32
	DelayBlkIOTicks      uint64
	GuestTime            uint64
	ChildGuestTime       int64
	StartData            uint64
	EndData              uint64
	StartBRK             uint64
	ArgStart             uint64
	ArgEnd               uint64
	EnvStart             uint64
	EnvEnd               uint64
	ExitCode             int32
}

// ParseStatString parse a string that contains the full contents of a
// Linux /proc/[pid]/stat file and returns the result in the first
// return value, and any error encountered calling fmt.Sscanf() on the
// string in the error return value.
func ParseStatString(s string) (Stat, error) {
	var ret Stat

	i := strings.LastIndex(s, ")")
	if i == -1 {
		return ret, errors.New("expected ')'")
	}

	if len(s) < i+2 {
		return ret, errors.New("input string too short")
	}

	_, err := fmt.Sscanf(strings.TrimSpace(s[i+1:]), "%c %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n",
		&ret.State,
		&ret.PPID,
		&ret.PGID,
		&ret.SessionID,
		&ret.TTYNumber,
		&ret.TPGID,
		&ret.Flags,
		&ret.MinorFaults,
		&ret.ChildMinorFaults,
		&ret.MajorFaults,
		&ret.ChildMajorFaults,
		&ret.UserTime,
		&ret.SystemTime,
		&ret.ChildUserTime,
		&ret.ChildSystemTime,
		&ret.Priority,
		&ret.Nice,
		&ret.NumThreads,
		&ret.ITRealValue,
		&ret.StartTime,
		&ret.VirtualMemSize,
		&ret.ResidentSetSize,
		&ret.ResidentSetSizeLimit,
		&ret.StartCode,
		&ret.EndCode,
		&ret.StartStack,
		&ret.KStackESP,
		&ret.KStackEIP,
		&ret.Signal,
		&ret.Blocked,
		&ret.SigIgnore,
		&ret.SigCatch,
		&ret.WChan,
		&ret.NumSwap,
		&ret.CumNumSwap,
		&ret.ExitSignal,
		&ret.Processor,
		&ret.RTPrio,
		&ret.Policy,
		&ret.DelayBlkIOTicks,
		&ret.GuestTime,
		&ret.ChildGuestTime,
		&ret.StartData,
		&ret.EndData,
		&ret.StartBRK,
		&ret.ArgStart,
		&ret.ArgEnd,
		&ret.EnvStart,
		&ret.EnvEnd,
		&ret.ExitCode)

	if err != nil {
		return ret, err
	}

	_, err = fmt.Sscanf(s[0:i+1], "%d", &ret.PID)
	if err != nil {
		return ret, err
	}

	h := strings.Index(s, "(")
	if h == -1 {
		return ret, errors.New("expected '('")
	}

	if !(h < i) {
		return ret, errors.New("expected '(' to come before ')'")
	}

	ret.Comm = s[h+1 : i]

	return ret, err
}

// ParseStatBytes is like ParseStatString but it operates on a bytes
// slice which contains the full contents of a /proc/[pid]/stat file.
func ParseStatBytes(b []byte) (Stat, error) {
	s := string(b)

	return ParseStatString(s)
}

// ParseStatReader consumes all the byes from r and calls ParseStatBytes on the result.
func ParseStatReader(r io.Reader) (Stat, error) {
	var ret Stat

	b, err := io.ReadAll(r)
	if err != nil {
		return ret, err
	}

	return ParseStatBytes(b)
}

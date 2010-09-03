// Copyright 2010 GoDCCP Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

// DCCP is implemented after http://read.cs.ucla.edu/dccp/rfc4340.txt

package dccp

import (
	"os"
)

// Network byte order (most significant byte first)

const (
	Modulo = 1 << 48	// 2^48
)

// Circular-arithmetic
//
// "It may make sense to store DCCP
// sequence numbers in the most significant 48 bits of 64-bit integers
// and set the least significant 16 bits to zero, since this supports a
// common technique that implements circular comparison A < B by testing
// whether (A - B) < 0 using conventional two's-complement arithmetic."

// "Reserved bitfields in DCCP packet headers MUST be set to zero by
// senders and MUST be ignored by receivers, unless otherwise specified."

// Half-connections

// NOTATION: F/X = feature number and endpoint; in F/A: 
// 'feature location'=A, 'feature remote'=B

// Default round-trip time for use when no estimate is available
// This parameter should default to not less than 0.2 seconds.
const DefaultRoundtripTime = ?

// The maximum segment lifetime, or MSL, is the maximum length of time a
//  packet can survive in the network.  The DCCP MSL should equal that of
//  TCP, which is normally two minutes.
const MaximumSegmentLifetime = ?

// Connections progress through three phases: initiation, including a three-way
// handshake; data transfer; and termination.

// DCCP sequence numbers increment by one per packet
type SequenceNumber uint64

// The nine possible states are as follows.  They are listed in
// increasing order.
const (
	CLOSED   = iota
	LISTEN   = _
	REQUEST  = _
	RESPOND  = _
	PARTOPEN = _
	OPEN     = _
	CLOSEREQ = _
	CLOSING  = _
	TIMEWAIT = _
)

// Congestion control mechanisms are denoted by one-byte 
// congestion control identifiers, or CCIDs.
// CCIDs 2 (TCP-like) and 3 (TFRC) are currently defined.

// There are four feature negotiation options in all: 
// Change L, Confirm L, Change R, and Confirm R.
// L = feature location, R = feature remote


// Up to Section 6 (non-inclusive)
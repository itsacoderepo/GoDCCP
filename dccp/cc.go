// Copyright 2011-2013 GoDCCP Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package dccp

// Regarding options and Half-Connection CCIDs (from Section 10.3):
//
// Any packet may contain information meant for either half-connection,
// so CCID-specific option types, feature numbers, and Reset Codes
// explicitly signal the half-connection to which they apply.
//
// o  Option numbers 128 through 191 are for options sent from the
//    HC-Sender to the HC-Receiver; option numbers 192 through 255 are
//    for options sent from the HC-Receiver to the HC-Sender.
//
// o  Reset Codes 128 through 191 indicate that the HC-Sender reset the
//    connection (most likely because of some problem with
//    acknowledgements sent by the HC-Receiver).  Reset Codes 192
//    through 255 indicate that the HC-Receiver reset the connection
//    (most likely because of some problem with data packets sent by the
//    HC-Sender).
//
// o  Finally, feature numbers 128 through 191 are used for features
//    located at the HC-Sender; feature numbers 192 through 255 are for
//    features located at the HC-Receiver.  Since Change L and Confirm L
//    options for a feature are sent by the feature location, we know
//    that any Change L(128) option was sent by the HC-Sender, while any
//    Change L(192) option was sent by the HC-Receiver.  Similarly,
//    Change R(128) options are sent by the HC-Receiver, while Change
//    R(192) options are sent by the HC-Sender.

// SenderCongestionControl specifies the interface for the congestion control logic of a DCCP
// sender (aka Half-Connection Sender CCID)
type SenderCongestionControl interface {

	// The congestion control is considered active between a call to Open and a call to Close or
	// an internal event that forces closure (like a reset event).

	// GetID() returns the CCID of this congestion control algorithm
	GetID() byte

	// GetCCMPS returns the Congestion Control Maximum Packet Size, CCMPS. Generally, PMTU <= CCMPS
	GetCCMPS() int32

	// GetRTT returns the Round-Trip Time as measured by this CCID
	GetRTT() int64

	// Open tells the Congestion Control that the connection has entered
	// OPEN or PARTOPEN state and that the CC can now kick in. Before the
	// call to Open and after the call to Close, the Strobe function is
	// expected to return immediately.
	Open()

	// Conn calls OnWrite before a packet is sent to give CongestionControl
	// an opportunity to add CCVal and options to an outgoing packet
	// NOTE: If the CC is not active, OnWrite should return 0, nil.
	OnWrite(ph *PreHeader) (ccval int8, options []*Option)

	// Conn calls OnRead after a packet has been accepted and validated
	// If OnRead returns ErrDrop, the packet will be dropped and no further processing
	// will occur. If OnRead returns ResetError, the connection will be reset.
	// NOTE: If the CC is not active, OnRead MUST return nil.
	OnRead(fb *FeedbackHeader) error

	// Strobe blocks until a new packet can be sent without violating the
	// congestion control rate limit. 
	// NOTE: If the CC is not active, Strobe MUST return immediately.
	Strobe()

	// OnIdle is called periodically, giving the CC a chance to:
	// (a) Request a connection reset by returning a CongestionReset, or
	// (b) Request the injection of an Ack packet by returning a CongestionAck
	// NOTE: If the CC is not active, OnIdle MUST to return nil.
	OnIdle(now int64) error

	// SetHeartbeat advices the CCID of the desired frequency of heartbeat packets.
	// A heartbeat interval value of zero indicates that no heartbeat is needed.
	SetHeartbeat(interval int64)

	// Close terminates the half-connection congestion control when it is not needed any longer
	Close()
}

// ReceiverCongestionControl specifies the interface for the congestion control logic of a DCCP
// receiver (aka Half-Connection Receiver CCID)
type ReceiverCongestionControl interface {

	// GetID() returns the CCID of this congestion control algorithm
	GetID() byte

	// Open tells the Congestion Control that the connection has entered
	// OPEN or PARTOPEN state and that the CC can now kick in.
	Open()

	// Conn calls OnWrite before a packet is sent to give CongestionControl
	// an opportunity to add CCVal and options to an outgoing packet
	// NOTE: If the CC is not active, OnWrite MUST return nil.
	OnWrite(ph *PreHeader) (options []*Option)

	// Conn calls OnRead after a packet has been accepted and validated
	// If OnRead returns ErrDrop, the packet will be dropped and no further processing
	// will occur. 
	// NOTE: If the CC is not active, OnRead MUST return nil.
	OnRead(ff *FeedforwardHeader) error

	// OnIdle behaves identically to the same method of the HC-Sender CCID
	OnIdle(now int64) error

	// Close terminates the half-connection congestion control when it is not needed any longer
	Close()
}

// PreHeader contains information that is shown to the 
// sender and receiver congesion controls before a packet is sent.
// PreHeader contains the parts of the DCCP header than are fixed before the
// CCID has made its changes to CCVal and Options.
type PreHeader struct {
	Type  byte
	X     bool
	SeqNo int64
	AckNo int64

	// TimeInject is the time when the packet was injected into the write
	// queue. This is either in the readLoop in response to a received
	// packet, in the idleLoop in response to idleness, or in the user
	// facing Write method. TimeInject is currently commented out,
	// since it is not used by the CC logic.
	// TimeInject int64

	// TimeWrite is the time just before the packet is handed off to the
	// network layer.  The bulk of the time difference between TimeWrite and
	// TimeMake reflects the duration that the packet spends in the write
	// queue before it is pulled for sending, which itself is a function of
	// the send rate. TimeWrite is currently used by the CC logic to measure
	// the roundtrip time without factoring rate-related wait times in
	// endpoint queues.
	TimeWrite int64
}

// FeedbackHeader contains information that is shown to the 
// sender congesion control at packet reception.
// FeedbackHeader encloses the parts of the packet header that
// are sent by the HC-Receiver and received by the HC-Sender
type FeedbackHeader struct {
	Type    byte
	X       bool
	SeqNo   int64
	Options []*Option
	AckNo   int64

	// Time when header received
	Time    int64
}

// FeedforwardHeader contains information that is shown to the 
// receiver congesion control at packet reception.
// FeedforwardHeader encloses the parts of the packet header that
// are sent by the HC-Sender and received by the HC-Receiver
type FeedforwardHeader struct {

	// These fields are copied directly from the header
	Type    byte
	X       bool
	SeqNo   int64
	CCVal   int8
	Options []*Option

	// Time when header received
	Time int64

	// Length of application data in bytes
	DataLen int
}

// CCID is a factory type that creates instances of sender and receiver CCIDs
type CCID interface {
	NewSender(env *Env, amb *Amb, args ...interface{}) SenderCongestionControl
	NewReceiver(env *Env, amb *Amb, args ...interface{}) ReceiverCongestionControl
}

const (
	CCID2 = 2 // TCP-like Congestion Control, RFC 4341
	CCID3 = 3 // TCP-Friendly Rate Control (TFRC), RFC 4342
)

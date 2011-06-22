// Copyright 2010 GoDCCP Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package retransmit

// All sequence numbers are circular integers

type header struct {

	// --> Sync indicates if the header includes a sync request, whose number is SyncNo
	Sync bool

	// The presence of SyncNo indicates that the sender is requesting an ack.
	// SyncNo is the sequence number of the sync request, generated by the sender.
	SyncNo uint16

	// --> Ack indicates if the header includes an acknowledgement, represented by SyncAckN
	Ack bool

	// AckSyncNo is the SyncNo this acknowledgement is in response to
	AckSyncNo uint16
	
	// AckDataNo is the sequence number of the first block in the acknowledgements map.
	// This number also indicates that all data blocks with sequence numbers lower 
	// than AckDataNo have already been received.
	AckDataNo uint32

	// AckMap is a bitmap where 1's correspond to blocks that HAVE NOT been received
	// TODO: This can be made more space-efficient with something like run-length encoding
	// or a Huffman-type (LZ?) thing that favors 0's
	AckMap []byte

	// --> Data indicates if the header includes data, represented by DataNo and DataCargo
	Data bool

	// DataNo is the sequence Number of the data block, if len(Data) > 0
	DataNo uint32

	// DataCargo is the contents of the data block
	DataCargo []byte
}

// Header wire format
//
//     +------------+---------------+----------------+----------------+
//     | Type 1byte | Ack Subheader | Sync Subheader | Data Subheader |
//     +------------+---------------+----------------+----------------+
//
// Type-byte format
//
//     MSB           LSB
//     +-+-+-+-+-+-+-+-+
//     | | | | | |D|S|A|
//     +-+-+-+-+-+-+-+-+
//
// Ack Subheader wire format
//
//     +------------------+------------------+------------------+----------------+
//     | AckSyncNo 2bytes | AckDataNo 4bytes | AckMapLen 2bytes | ... AckMap ... |
//     +------------------+------------------+------------------+----------------+
//
// Sync Subheader wire format
//
//     +----------------+
//     | SyncNo 2 bytes |
//     +----------------+
//
// Data Subheader wire format
//
//     +---------------+----------------+--------------+
//     | DataNo 4bytes | DataLen 2bytes | ... Data ... |
//     +---------------+----------------+--------------+

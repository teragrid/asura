# Asura - Application BlockChain Interface (asura)

[![CircleCI](https://circleci.com/gh/teragrid/asura.svg?style=svg)](https://circleci.com/gh/teragrid/asura)

Blockchains are systems for multi-master state machine replication.
**Asura** is an interface that defines the boundary between the replication engine (the blockchain),
and the state machine (the application).
By using a socket protocol, we enable a consensus engine running in one process
to manage an application state running in another.

For background information on asura, motivations, and teragrid, please visit [the documentation](http://teragrid.readthedocs.io/en/master/).
The two guides to focus on are the `Application Development Guide` and `Using asura-CLI`.

Previously, the asura was referred to as TMSP.

The community has provided a number of addtional implementations, see the [Teragrid Ecosystem](https://teragrid.network/ecosystem)

## Specification

The [primary specification](https://github.com/teragrid/asura/blob/master/types/types.proto)
is made using Protocol Buffers. To build it, run

```
make protoc
```

See `protoc --help` and [the Protocol Buffers site](https://developers.google.com/protocol-buffers)
for details on compiling for other languages. Note we also include a [GRPC](http://www.grpc.io/docs)
service definition.

For the specification as an interface in Go, see the
[types/application.go file](https://github.com/teragrid/asura/blob/master/types/application.go).

See the [spec file](specification.rst) for a detailed description of the message types.

## Install

```
go get github.com/teragrid/asura
cd $GOPATH/src/github.com/teragrid/asura
make get_vendor_deps
make install
```

## Implementation

We provide three implementations of the asura in Go:

- Golang in-process
- asura-socket
- GRPC

Note the GRPC version is maintained primarily to simplify onboarding and prototyping and is not receiving the same
attention to security and performance as the others

### In Process

The simplest implementation just uses function calls within Go.
This means asura applications written in Golang can be compiled with TeragridCore and run as a single binary.

See the [examples](#examples) below for more information.

### Socket (TSP)

Asura is best implemented as a streaming protocol.
The socket implementation provides for asynchronous, ordered message passing over unix or tcp.
Messages are serialized using Protobuf3 and length-prefixed with a [signed Varint](https://developers.google.com/protocol-buffers/docs/encoding?csw=1#signed-integers)

For example, if the Protobuf3 encoded asura message is `0xDEADBEEF` (4 bytes), the length-prefixed message is `0x08DEADBEEF`, since `0x08` is the signed varint
encoding of `4`. If the Protobuf3 encoded asura message is 65535 bytes long, the length-prefixed message would be like `0xFEFF07...`.

Note the benefit of using this `varint` encoding over the old version (where integers were encoded as `<len of len><big endian len>` is that
it is the standard way to encode integers in Protobuf. It is also generally shorter.

### GRPC

GRPC is an rpc framework native to Protocol Buffers with support in many languages.
Implementing the asura using GRPC can allow for faster prototyping, but is expected to be much slower than
the ordered, asynchronous socket protocol. The implementation has also not received as much testing or review.

Note the length-prefixing used in the socket implementation does not apply for GRPC.

## Usage

The `asura-cli` tool wraps an asura client and can be used for probing/testing an asura server.
For instance, `asura-cli test` will run a test sequence against a listening server running the Counter application (see below).
It can also be used to run some example applications.
See [the documentation](http://teragrid.readthedocs.io/en/master/) for more details.

### Examples

Check out the variety of example applications in the [example directory](example/).
It also contains the code refered to by the `counter` and `kvstore` apps; these apps come
built into the `asura-cli` binary.

#### Counter

The `asura-cli counter` application illustrates nonce checking in transactions. It's code looks like:

```golang
func cmdCounter(cmd *cobra.Command, args []string) error {

	app := counter.NewCounterApplication(flagSerial)

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	// Start the listener
	srv, err := server.NewServer(flagAddrC, flagasura, app)
	if err != nil {
		return err
	}
	srv.SetLogger(logger.With("module", "asura-server"))
	if err := srv.Start(); err != nil {
		return err
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
	return nil
}
```

and can be found in [this file](cmd/asura-cli/asura-cli.go).

#### kvstore

The `asura-cli kvstore` application, which illustrates a simple key-value Merkle tree

```golang
func cmdKVStore(cmd *cobra.Command, args []string) error {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	// Create the application - in memory or persisted to disk
	var app types.Application
	if flagPersist == "" {
		app = kvstore.NewKVStoreApplication()
	} else {
		app = kvstore.NewPersistentKVStoreApplication(flagPersist)
		app.(*kvstore.PersistentKVStoreApplication).SetLogger(logger.With("module", "kvstore"))
	}

	// Start the listener
	srv, err := server.NewServer(flagAddrD, flagasura, app)
	if err != nil {
		return err
	}
	srv.SetLogger(logger.With("module", "asura-server"))
	if err := srv.Start(); err != nil {
		return err
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
	return nil
}
```

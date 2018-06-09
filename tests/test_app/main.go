package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/teragrid/asura/example/code"
	"github.com/teragrid/asura/types"
)

var asuraType string

func init() {
	asuraType = os.Getenv("asura")
	if asuraType == "" {
		asuraType = "socket"
	}
}

func main() {
	testCounter()
}

const (
	maxasuraConnectTries = 10
)

func ensureasuraIsUp(typ string, n int) error {
	var err error
	cmdString := "asura-cli echo hello"
	if typ == "grpc" {
		cmdString = "asura-cli --asura grpc echo hello"
	}

	for i := 0; i < n; i++ {
		cmd := exec.Command("bash", "-c", cmdString) // nolint: gas
		_, err = cmd.CombinedOutput()
		if err == nil {
			break
		}
		<-time.After(500 * time.Millisecond)
	}
	return err
}

func testCounter() {
	asuraApp := os.Getenv("asura_APP")
	if asuraApp == "" {
		panic("No asura_APP specified")
	}

	fmt.Printf("Running %s test with asura=%s\n", asuraApp, asuraType)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("asura-cli %s", asuraApp)) // nolint: gas
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		log.Fatalf("starting %q err: %v", asuraApp, err)
	}
	defer cmd.Wait()
	defer cmd.Process.Kill()

	if err := ensureasuraIsUp(asuraType, maxasuraConnectTries); err != nil {
		log.Fatalf("echo failed: %v", err)
	}

	client := startClient(asuraType)
	defer client.Stop()

	setOption(client, "serial", "on")
	commit(client, nil)
	deliverTx(client, []byte("abc"), code.CodeTypeBadNonce, nil)
	commit(client, nil)
	deliverTx(client, []byte{0x00}, types.CodeTypeOK, nil)
	commit(client, []byte{0, 0, 0, 0, 0, 0, 0, 1})
	deliverTx(client, []byte{0x00}, code.CodeTypeBadNonce, nil)
	deliverTx(client, []byte{0x01}, types.CodeTypeOK, nil)
	deliverTx(client, []byte{0x00, 0x02}, types.CodeTypeOK, nil)
	deliverTx(client, []byte{0x00, 0x03}, types.CodeTypeOK, nil)
	deliverTx(client, []byte{0x00, 0x00, 0x04}, types.CodeTypeOK, nil)
	deliverTx(client, []byte{0x00, 0x00, 0x06}, code.CodeTypeBadNonce, nil)
	commit(client, []byte{0, 0, 0, 0, 0, 0, 0, 5})
}

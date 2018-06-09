package testsuite

import (
	"bytes"
	"errors"
	"fmt"

	asuracli "github.com/teragrid/asura/client"
	"github.com/teragrid/asura/types"
	cmn "github.com/teragrid/teralibs/common"
)

func InitChain(client asuracli.Client) error {
	total := 10
	vals := make([]types.Validator, total)
	for i := 0; i < total; i++ {
		pubkey := cmn.RandBytes(33)
		power := cmn.RandInt()
		vals[i] = types.Validator{pubkey, int64(power)}
	}
	_, err := client.InitChainSync(types.RequestInitChain{
		Validators:    vals,
		AppStateBytes: []byte("{}"),
	})
	if err != nil {
		fmt.Printf("Failed test: InitChain - %v\n", err)
		return err
	}
	fmt.Println("Passed test: InitChain")
	return nil
}

func SetOption(client asuracli.Client, key, value string) error {
	_, err := client.SetOptionSync(types.RequestSetOption{Key: key, Value: value})
	if err != nil {
		fmt.Println("Failed test: SetOption")
		fmt.Printf("error while setting %v=%v: \nerror: %v\n", key, value, err)
		return err
	}
	fmt.Println("Passed test: SetOption")
	return nil
}

func Commit(client asuracli.Client, hashExp []byte) error {
	res, err := client.CommitSync()
	data := res.Data
	if err != nil {
		fmt.Println("Failed test: Commit")
		fmt.Printf("error while committing: %v\n", err)
		return err
	}
	if !bytes.Equal(data, hashExp) {
		fmt.Println("Failed test: Commit")
		fmt.Printf("Commit hash was unexpected. Got %X expected %X\n", data, hashExp)
		return errors.New("CommitTx failed")
	}
	fmt.Println("Passed test: Commit")
	return nil
}

func DeliverTx(client asuracli.Client, txBytes []byte, codeExp uint32, dataExp []byte) error {
	res, _ := client.DeliverTxSync(txBytes)
	code, data, log := res.Code, res.Data, res.Log
	if code != codeExp {
		fmt.Println("Failed test: DeliverTx")
		fmt.Printf("DeliverTx response code was unexpected. Got %v expected %v. Log: %v\n",
			code, codeExp, log)
		return errors.New("DeliverTx error")
	}
	if !bytes.Equal(data, dataExp) {
		fmt.Println("Failed test: DeliverTx")
		fmt.Printf("DeliverTx response data was unexpected. Got %X expected %X\n",
			data, dataExp)
		return errors.New("DeliverTx error")
	}
	fmt.Println("Passed test: DeliverTx")
	return nil
}

func CheckTx(client asuracli.Client, txBytes []byte, codeExp uint32, dataExp []byte) error {
	res, _ := client.CheckTxSync(txBytes)
	code, data, log := res.Code, res.Data, res.Log
	if code != codeExp {
		fmt.Println("Failed test: CheckTx")
		fmt.Printf("CheckTx response code was unexpected. Got %v expected %v. Log: %v\n",
			code, codeExp, log)
		return errors.New("CheckTx")
	}
	if !bytes.Equal(data, dataExp) {
		fmt.Println("Failed test: CheckTx")
		fmt.Printf("CheckTx response data was unexpected. Got %X expected %X\n",
			data, dataExp)
		return errors.New("CheckTx")
	}
	fmt.Println("Passed test: CheckTx")
	return nil
}

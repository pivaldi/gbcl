package main

import (
	"fmt"
	"os"
	"time"

	"piprim.net/gbcl/app/db"
	dbblock "piprim.net/gbcl/app/db/block"
	"piprim.net/gbcl/app/tx"
	apptype "piprim.net/gbcl/app/type"
	liberrors "piprim.net/gbcl/lib/errors"
)

//nolint:gomnd // Because one shot script
func main() {
	state, err := db.NewStateFromDisk()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer state.Close()

	block0 := dbblock.New(
		apptype.Hash{},
		uint64(time.Now().Unix()),
		[]tx.Tx{
			tx.New("root", "root", 3, ""),
			tx.New("root", "root", 700, "reward"),
		},
	)

	err = state.AddBlock(block0)
	if err != nil {
		liberrors.HandleError(err)
	}

	block0hash, _ := state.Persist()

	block1 := dbblock.New(
		block0hash,
		uint64(time.Now().Unix()),
		[]tx.Tx{
			tx.New("root", "babayaga", 2000, ""),
			tx.New("root", "root", 100, "reward"),
			tx.New("babayaga", "root", 1, ""),
			tx.New("babayaga", "caesar", 1000, ""),
			tx.New("babayaga", "root", 50, ""),
			tx.New("root", "root", 600, "reward"),
		},
	)

	err = state.AddBlock(block1)
	if err != nil {
		liberrors.HandleError(err)
	}

	_, err = state.Persist()
	if err != nil {
		liberrors.HandleError(err)
	}
}

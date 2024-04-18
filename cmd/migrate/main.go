package main

import (
	"time"

	"piprim.net/gbcl/app"
	"piprim.net/gbcl/app/db"
	dbblock "piprim.net/gbcl/app/db/block"
	"piprim.net/gbcl/app/tx"
	apptype "piprim.net/gbcl/app/type"
)

//nolint:gomnd // Because one shot script
func main() {
	err := app.Init("")
	if err != nil {
		panic(err)
	}

	state, err := db.NewStateFromDisk()
	if err != nil {
		panic(err)
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
		panic(err)
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
		panic(err)
	}

	_, err = state.Persist()
	if err != nil {
		panic(err)
	}
}

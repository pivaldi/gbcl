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
		state.NextBlockNumber(),
		uint64(time.Now().Unix()),
		[]tx.Tx{
			tx.New("root", "root", 3, ""),
			tx.New("root", "root", 700, "reward"),
		},
	)

	hash0, err := state.AddBlock(block0)
	if err != nil {
		panic(err)
	}

	block1 := dbblock.New(
		*hash0,
		state.NextBlockNumber(),
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

	hash1, err := state.AddBlock(block1)
	if err != nil {
		panic(err)
	}

	block2 := dbblock.New(
		*hash1,
		state.NextBlockNumber(),
		uint64(time.Now().Unix()),
		[]tx.Tx{
			tx.New("root", "root", 24700, "reward"),
		},
	)

	_, err = state.AddBlock(block2)
	if err != nil {
		panic(err)
	}
}

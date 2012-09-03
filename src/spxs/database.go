package main

import (
	"bufio"
	"io"
	"log"
	"os"
	. "spexs"
	"strings"
)

func CreateDatabase(conf Conf) *Database {
	db := NewDatabase(1024)

	if conf.Alphabet.Characters == "" {
		log.Fatal("No alphabet defined!")
	}

	for _, alpha := range conf.Alphabet.Characters {
		db.AddToken(string(alpha))
	}

	for id, grp := range conf.Alphabet.Groups {
		group := Group{}

		if len(id) != 1 {
			log.Fatal("Group identifier must be of length 1.")
		}

		group.Alias = id
		group.Str = "[" + grp.Group + "]"

		tokens := strings.Split(grp.Group, "")
		group.Elems = db.ToTids(tokens)

		db.AddGroup(group)
	}

	if conf.Data.Input == "" {
		log.Fatal("Data file not defined")
	}

	addSeqsFromFile(db, conf.Data.Input, 0)

	if conf.Data.Database != "" {
		addSeqsFromFile(db, conf.Data.Database, 1)
	}

	return db
}

func addSeqsFromFile(db *Database, filename string, section int) {
	var (
		file   *os.File
		reader *bufio.Reader
		line   string
		err    error
	)

	if file, err = os.Open(filename); err != nil {
		log.Println("Did not find data file:", filename)
		log.Fatal(err)
	}

	reader = bufio.NewReader(file)
	for err == nil {
		if line, err = reader.ReadString('\n'); err != nil && err != io.EOF {
			log.Fatal(err)
		}

		line = strings.TrimSpace(line)
		tokens := strings.Split(line, "")
		tids := db.ToTids(tokens)

		if len(tids) <= 0 {
			continue
		}

		seq := Sequence{tids, len(tids), section, 1}
		db.AddSequence(seq)
	}
}
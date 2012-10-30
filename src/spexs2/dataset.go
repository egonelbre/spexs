package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path"
	. "spexs"
	"strconv"
	"strings"
)

type File struct {
	Id   int
	Base string
	Full string
}

type Dataset struct {
	Groups map[string][]int
	Files  map[string]File
}

func NewDataset() *Dataset {
	return &Dataset{make(map[string][]int), make(map[string]File)}
}

func CreateDatabase(conf *Conf) (*Database, *Dataset) {
	db := NewDatabase(1024)
	db.Separator = conf.Alphabet.Separator

	ds := NewDataset()

	for id, grp := range conf.Alphabet.Groups {
		group := TokenGroup{}

		group.Name = id
		group.FullName = "[" + grp.Elements + "]"

		tokenNames := strings.Split(grp.Elements, db.Separator)
		tokenNames = removeEmpty(tokenNames)
		group.Elems = db.ToTokens(tokenNames)

		db.AddGroup(group)
	}

	ds.AddFileGroups(db, conf.Data)
	return db, ds
}

func loadFileList(filename string) []string {
	var (
		file  *os.File
		err   error
		line  string
		lines []string
	)

	lines = make([]string, 0)
	if file, err = os.Open(filename); err != nil {
		log.Println("Did not find file list:", filename)
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)
	for err == nil {
		if line, err = reader.ReadString('\n'); err != nil && err != io.EOF {
			log.Fatal(err)
		}

		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}

	}
	return lines
}

func (ds *Dataset) AddFileGroups(db *Database, groups map[string]FileGroup) {
	for group, filegroup := range groups {

		files := make([]string, 0)
		if filegroup.File != "" {
			files = append(files, filegroup.File)
		}
		files = append(files, filegroup.Files...)

		if filegroup.FileList != "" {
			files = append(files, loadFileList(filegroup.FileList)...)
		}

		ids := make([]int, 0)
		for _, file := range files {
			id := ds.AddFile(db, file, filegroup.CountSeparator)
			ids = append(ids, id)
		}
		ds.Groups[group] = ids
	}
}

func removeEmpty(names []string) []string {
	result := make([]string, len(names))
	i := 0
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed != "" {
			result[i] = trimmed
			i += 1
		}
	}
	return result[0:i]
}

func (ds *Dataset) AddFile(db *Database, filename string, countSeparator string) int {
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

	isCounted := countSeparator != ""

	section := db.MakeSection()
	reader = bufio.NewReader(file)
	for err == nil {
		count := 1

		if isCounted {
			if line, err = reader.ReadString(countSeparator[1]); err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if count, err = strconv.Atoi(line); err != nil {
				log.Fatal(err)
			}
		}

		if line, err = reader.ReadString('\n'); err != nil && err != io.EOF {
			log.Fatal(err)
		}

		line = strings.TrimSpace(line)
		tokenNames := strings.Split(line, db.Separator)
		tokenNames = removeEmpty(tokenNames)

		if len(tokenNames) <= 0 {
			continue
		}

		db.AddSequence(section, tokenNames, count)
	}

	name := path.Base(filename)
	ds.Files[name] = File{section, name, filename}
	return section
}

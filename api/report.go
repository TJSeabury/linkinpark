package main

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
)

func creatReport(root map[string]pageInfo, filename string) []byte {
	keys := make([]string, 0, len(root))
	for k := range root {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	csvFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	var _p pageInfo
	csvwriter.Write(_p.writeCSVLabels())
	for _, k := range keys {
		p := root[k]
		csvwriter.Write(p.writeCSVLine())
	}
	csvwriter.Flush()

	data, err := os.ReadFile(filename)
	check(err)

	csvFile.Close()

	//log.Println("Deleteing", filename, " . . . ")
	err = os.Remove(filename)
	check(err)

	return data
}

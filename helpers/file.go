package helpers

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

func ReadCsv(filepath string) [][]string {
	f, err := os.Open(filepath)

	if err != nil {
		log.Fatalln("Unable to read input path "+filepath, err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	record, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(`Unable to parse file as csv from `+filepath, err)
	}
	return record
}

func ReadJSON[V interface{}](filepath string) V {
	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}
	var data V
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

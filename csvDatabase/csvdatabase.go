package csvdatabase

import (
	"bytes"
	"encoding/csv"
	"fmt"
	datastorage "github/bustman/shops/dataStorage"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type csvDB struct {
	dataCsv [][]string
}

// if not exist creat a csv file with desired FileName, if FileName allready exist delete it and create New One
func CreatUpdateCsvFile(FileName string, c csvDB) error {
	if _, err := os.Stat(FileName + ".csv"); err == nil {
		err := os.Remove(FileName + ".csv")
		if err != nil {
			log.Fatal(err)
		}
	}
	f, err := os.OpenFile(FileName+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		f.Close()
		return err
	}
	w := csv.NewWriter(f)
	w.WriteAll(c.dataCsv)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return nil
}
func ReadCsvFile(FileName string) (*csvDB, error) {
	f, err := os.ReadFile(FileName + ".csv")
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(bytes.NewReader(f))
	rec, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return &csvDB{dataCsv: rec}, nil
}

func (cData *CsvData) PutToCsv(name string) error {
	var csvDB csvDB
	data := cData.data[name]
	if values, ok := data.(map[string]float64); !ok {
		log.Fatal("PutToCsv func Invalide csv type")
	} else {
		for k, v := range values {
			var list []string
			list = append(list, k, strconv.FormatFloat(v, 'G', -1, 32))
			csvDB.dataCsv = append(csvDB.dataCsv, list)
		}
	}
	err := CreatUpdateCsvFile(name, csvDB)
	if err != nil {
		return err
	}
	return nil
}

type CsvData struct {
	data map[string]interface{}
}

func (cData *CsvData) GetFromCsv(FileName string) error {
	if cData.data == nil {
		cData.data = make(map[string]interface{})
	}
	csvDB, err := ReadCsvFile(FileName)
	if err != nil {
		return err
	}
	m := make(map[string]float64)
	for _, v := range csvDB.dataCsv {
		if len(v) != 2 {
			log.Fatal("colone data error")
		}
		n, err := strconv.ParseFloat(v[1], 32)
		if err != nil {
			log.Fatal(err)
		}
		m[v[0]] = n
	}
	cData.data[FileName] = m
	return nil
}

func (csv *CsvData) Load(FileName string) {
	if err := csv.GetFromCsv(FileName); err != nil {
		log.Fatal(err)
	}

}

// Save data in csv file with the given FileName
func (csv *CsvData) Save(FileName string, dataValue interface{}) {
	if csv.data == nil {
		csv.data = make(map[string]interface{})
	}
	if values, ok := dataValue.(map[string]float64); !ok {
		log.Fatal("csv.Save Invalide csv type")
	} else {
		csv.data[FileName] = values
		if err := csv.PutToCsv(FileName); err != nil {
			log.Fatal("err csv.PutToCsv", err)
		}
	}
	csv.data = nil
}
func (csv *CsvData) List() string {
	if csv.data == nil {
		csv.LoadAll()
	}
	return csv.String()
}
func (csv *CsvData) String() string {
	var result string
	for key, db := range csv.data {
		result += fmt.Sprintf("Name : %s\n", key)
		result += "Data : shop | price\n"
		if data, ok := db.(map[string]float64); !ok {
			log.Fatal("Stringer type not ok")
		} else {
			for k, v := range data {
				result += fmt.Sprintf("      %s    %v\n", k, v)
			}
		}
	}
	return result
}
func GetFileList(dir string) ([]string, error) {
	var listFiles []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			listFiles = append(listFiles, strings.Trim(filepath.Base(path), ".csv"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return listFiles, nil
}
func (csv *CsvData) LoadAll() error {
	list, err := GetFileList(".")
	if err != nil {
		return err
	}
	for _, FileName := range list {
		csv.Load(FileName)
	}
	return nil
}
func NewCsvData() datastorage.Datastorage {
	return &CsvData{data: map[string]interface{}{}}
}

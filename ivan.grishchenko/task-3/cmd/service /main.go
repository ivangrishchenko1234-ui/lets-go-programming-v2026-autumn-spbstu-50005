package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type OutputValute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to configuration file")
	flag.Parse()

	configFile, err := os.ReadFile(*configPath)
	if err != nil {
		panic("failed to read config file: " + err.Error())
	}

	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		panic("failed to parse yaml config: " + err.Error())
	}

	xmlData, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		panic("failed to read input xml file: " + err.Error())
	}

	var valCurs ValCurs
	err = xml.Unmarshal(xmlData, &valCurs)
	if err != nil {
		panic("failed to parse xml data: " + err.Error())
	}

	var outputList []OutputValute
	for _, v := range valCurs.Valutes {
		valStr := strings.ReplaceAll(v.Value, ",", ".")
		parsedVal, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			panic("failed to parse currency value: " + err.Error())
		}

		outputList = append(outputList, OutputValute{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    parsedVal,
		})
	}

	sort.Slice(outputList, func(i, j int) bool {
		return outputList[i].Value > outputList[j].Value
	})

	outputDir := filepath.Dir(cfg.OutputFile)
	err = os.MkdirAll(outputDir, 0755) 
	if err != nil {
		panic("failed to create output directory: " + err.Error())
	}

	jsonData, err := json.MarshalIndent(outputList, "", "  ")
	if err != nil {
		panic("failed to marshal json: " + err.Error())
	}

	err = os.WriteFile(cfg.OutputFile, jsonData, 0644)
	if err != nil {
		panic("failed to write output file: " + err.Error())
	}

	println("Successfully processed and saved result to:", cfg.OutputFile)
}

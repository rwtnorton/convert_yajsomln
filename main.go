package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

//const input = `
//a: Easy!
//b:
//  c: 2
//  d: [3, 4]
//`

var fromJSON bool
var fromYAML bool

func init() {
	flag.BoolVar(&fromJSON, "from_json", false, "treat stdin as JSON, output YAML")
	flag.BoolVar(&fromYAML, "from_yaml", false, "treat stdin as YAML, output JSON")
	flag.Parse()

	if fromJSON == fromYAML {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "must select only one of --from_json or --from_yaml\n")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	bs, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	s := string(bs)

	if fromYAML {
		jsonStr, err := yamlToJSON(s)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(jsonStr)
	}

	if fromJSON {
		yamlStr, err := jsonToYAML(s)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(yamlStr)
	}
}

func yamlToJSON(s string) (string, error) {
	var err error
	m := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(s), &m)
	if err != nil {
		return "", fmt.Errorf("unable to yaml.Unmarshal(): %q: %v", s, err)
	}

	bs, err := json.MarshalIndent(m, ``, `    `)
	if err != nil {
		return "", fmt.Errorf("unable to json.Marshal(): %+v: %v", m, err)
	}

	return string(bs), nil
}

func jsonToYAML(s string) (string, error) {
	var err error
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return "", fmt.Errorf("unable to json.Unmarshal(): %q: %v", s, err)
	}

	bs, err := yaml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("unable to yaml.Marshal(): %+v: %v", m, err)
	}

	return string(bs), nil
}

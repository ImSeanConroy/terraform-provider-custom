package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	definitionsPath  = "./gen/definitions/"
	providerTemplate = "./gen/templates/provider.go"
	providerLocation = "./internal/provider/provider.go"
)

type t struct {
	path   string
	prefix string
	suffix string
}

var templates = []t{
	// {
	// 	path:   "./gen/templates/data_source.go",
	// 	prefix: "./internal/provider/data_source_custom_",
	// 	suffix: ".go",
	// },
	// {
	// 	path:   "./gen/templates/data_source_test.go",
	// 	prefix: "./internal/provider/data_source_custom_",
	// 	suffix: "_test.go",
	// },
	// {
	// 	path:   "./gen/templates/resource.go",
	// 	prefix: "./internal/provider/resource_custom_",
	// 	suffix: ".go",
	// },
	// {
	// 	path:   "./gen/templates/resource_test.go",
	// 	prefix: "./internal/provider/resource_custom_",
	// 	suffix: "_test.go",
	// },
	{
		path:   "./gen/templates/data-source.tf",
		prefix: "./examples/data-sources/custom_",
		suffix: "/data-source.tf",
	},
	{
		path:   "./gen/templates/resource.tf",
		prefix: "./examples/resources/custom_",
		suffix: "/resource.tf",
	},
	{
		path:   "./gen/templates/import.sh",
		prefix: "./examples/resources/custom_",
		suffix: "/import.sh",
	},
}

type YamlConfig struct {
	Name           string                `yaml:"name"`
	RestEndpoint   string                `yaml:"rest_endpoint"`
	DsDescription  string                `yaml:"ds_description"`
	ResDescription string                `yaml:"res_description"`
	IdAttribute    string                `yaml:"id_attribute"`
	Attributes     []YamlConfigAttribute `yaml:"attributes"`
}

type YamlConfigAttribute struct {
	ModelName   string `yaml:"model_name"`
	TfName      string `yaml:"tf_name"`
	Type        string `yaml:"type"`
	Id          bool   `yaml:"id"`
	Reference   bool   `yaml:"reference"`
	Mandatory   bool   `yaml:"mandatory"`
	Description string `yaml:"description"`
	Example     string `yaml:"example"`
}

// Helper function to convert a string in snake case
func toSnakeCase(str string) string {
	return strings.ToLower(strings.ReplaceAll(str, " ", "_"))
}

var functions = template.FuncMap{
	"snakeCase": toSnakeCase,
}

func main() {
	// Get list of YAML files in the definition directory
	files, err := os.ReadDir(definitionsPath)
	if err != nil {
		log.Fatal("Error reading definitions directory: %w", err)
	}

	// Loop through YAML files and process each
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {

			// Read YAML file
			yamlFile, err := os.ReadFile(filepath.Join(definitionsPath, file.Name()))
			if err != nil {
				log.Fatalf("Error reading YAML file %s: %v", file.Name(), err)
				continue
			}

			// Parse YAML into struct
			config := YamlConfig{}
			err = yaml.Unmarshal(yamlFile, &config)
			if err != nil {
				log.Fatalf("Error parsing YAML file %s: %v", file.Name(), err)
				continue
			}

			for _, tmplInfo := range templates {
				// // Create a new template with the base name of the template file
				tmpl, err := template.New(path.Base(tmplInfo.path)).Funcs(functions).ParseFiles(tmplInfo.path)
				if err != nil {
					log.Printf("Error parsing template file %s: %v", tmplInfo.path, err)
					continue
				}

				// // Define output file path based on name field
				outputFileName := fmt.Sprintf("%s%s%s", tmplInfo.prefix, toSnakeCase(config.Name), tmplInfo.suffix)
				_, err = os.Open(outputFileName)
				if err != nil {
					os.MkdirAll(filepath.Dir(outputFileName), 0755)
				}

				// // // Create the output file
				outFile, err := os.Create(outputFileName)
				if err != nil {
					log.Fatalf("Error creating ouput file %s: %v", outputFileName, err)
					continue
				}
				defer outFile.Close()

				err = tmpl.Execute(outFile, config)
				if err != nil {
					log.Fatalf("Error executing template %s: %v", tmplInfo.path, err)
					continue
				}

			}
			fmt.Printf("Template rendering completed successfully. %s\n", file.Name())

		}
	}
}

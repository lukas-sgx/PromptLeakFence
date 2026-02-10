package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type filePolicy struct {
	Policy struct {
		Exclude []string `yaml:"exclude"`
	} `yaml:"policy"`
}

func ReadPolicy(file string) []string {
	policy := &filePolicy{}
	yamlFile, err := os.ReadFile(file)

	if err != nil {
		log.Printf("Error opening file: %s", file)
	}
	err = yaml.Unmarshal(yamlFile, policy)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return policy.Policy.Exclude
}

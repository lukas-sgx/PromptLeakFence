package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type FilePolicy struct {
	Policy struct {
		Exclude           []string `yaml:"exclude"`
		InjectionPatterns []string `yaml:"injection_patterns"`
	} `yaml:"policy"`
}

func ReadPolicy(file string) *FilePolicy {
	policy := &FilePolicy{}
	yamlFile, err := os.ReadFile(file)

	if err != nil {
		log.Printf("Error opening file: %s", file)
	}
	err = yaml.Unmarshal(yamlFile, policy)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return policy
}

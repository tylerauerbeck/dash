package inventory

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Inventory represents a collection of objects to be managed within an cluster
type Inventory struct {
	Version            int64                  `yaml:"version"`
	Namespace          string                 `yaml:"namespace"`
	ResourceGroups     []ResourceGroup        `yaml:"resource_groups"`
	ClusterContentList []ClusterContentObject `yaml:"openshift_cluster_content"`
	Prefix             string
}

// ResourceGroup represents a collection of objects within a V3 dash inventory
type ResourceGroup struct {
	Name      string     `yaml:"name"`
	Namespace string     `yaml:"namespace"`
	Resources []Resource `yaml:"resources"`
}

// Resource represents a single object within a ResourceGroup
type Resource struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	File      string `yaml:"file"`
	Action    Action `yaml:"action"`
}

// ClusterContentObject represents a single object within ClusterContentList
type ClusterContentObject struct {
	Object  string           `yaml:"object"`
	Content []ClusterContent `yaml:"content"`
}

// ClusterContent represents the actual content of a ClusterContentObject
type ClusterContent struct {
	Name           string `yaml:"name"`
	Namespace      string `yaml:"namespace,omitempty"`
	File           string `yaml:"file,omitempty"`
	Template       string `yaml:"template,omitempty"`
	Params         string `yaml:"params,omitempty"`
	ParamsFromVars string `yaml:"params_from_vars,omitempty"`
	Action         string `yaml:"action,omitempty"`
}

// Action is the type of action to be performed
type Action string

// UnmarshalJSON implements the Unmarshaler interface on Action, so we can default it to "apply"
func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*a = Action("apply")
	} else {
		*a = Action(s)
	}
	return nil
}

// Load performs the action of reading a dash.yml inventory
func (i *Inventory) Load(pre string) *Inventory {

	i.Prefix = pre

	yamlFile, err := ioutil.ReadFile(pre + "dash.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, i)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return i
}

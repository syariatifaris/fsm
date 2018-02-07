package sample

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type BusinessFSM struct {
	Fsm []ContainerFSM `json:"fsm"`
}

type ContainerFSM struct {
	State State `json:"state"`
}

type State struct {
	StateName   string       `json:"name"`
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	TargetState        string   `json:"target"`
	PostOperationHooks []string `json:"post_operation_hooks"`
}

func ReadFromFile(name string, v interface{}) error {
	var directories = []string{
		"files/fsm",
		"../files/fsm",
		"../../files/fsm",
		"../../../files/fsm",
		"/etc/fsm",
	}

	for _, dir := range directories {
		raw, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", dir, name))
		if err == nil {
			err := json.Unmarshal(raw, v)
			return err
		}
	}

	return errors.New(fmt.Sprintf("unable to resolve %s.json file", name))
}

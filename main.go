package main

import (
	"github.com/syariatifaris/fsm/kvfsm"
	"fmt"
)

func main(){
	fsm := kvfsm.NewFsm()
	fsm.AddFSM("initial", handleInitial, nil)
	fsm.AddFSM("first_transition", handleFirstTransition, nil)

	//get states
	first, err := fsm.GetState("first_transition")
	if err != nil{
		panic(err)
	}

	//set transition
	fsm.States["initial"].Transitions = make(map[string]kvfsm.Transition)
	fsm.States["initial"].Transitions["first_transition"] = kvfsm.Transition{
		State:first,
	}

	//add data dummy
	data := "a data being passed"

	//run change transition
	v, err := fsm.GoToStateFrom("initial", "first_transition", data)
	if err != nil{
		panic(err)
	}

	fmt.Println(fmt.Sprintf("processed %+v", v))
}

func handleInitial(fsm *kvfsm.TransactionState, v interface{}) (interface{}, error){
	fmt.Println(fmt.Sprintf("%s is being processed", fsm.StateName))
	fmt.Println(fmt.Sprintf("processing %+v", v))
	return v, nil
}

func handleFirstTransition(fsm *kvfsm.TransactionState, v interface{}) (interface{}, error){
	fmt.Println(fmt.Sprintf("%s is being processed", fsm.StateName))
	fmt.Println(fmt.Sprintf("processing %+v", v))
	return v, nil
}
package kvfsm

import (
	"errors"
	"fmt"
)

func NewFsm() *Fsm {
	return &Fsm{
		States: make(map[string]*TransactionState),
	}
}

type Fsm struct {
	States map[string]*TransactionState
}

type StateActionFunc func(*TransactionState, interface{}) (interface{}, error)

//arg interface from state action response
//arg error pass from state action response
type HookRequestFunc func([]string, interface{}, error) (interface{}, error)

type TransactionState struct {
	StateName   string
	StateAction StateActionFunc
	HookAction  HookRequestFunc
	Transitions map[string]Transition
}

type Transition struct {
	State              *TransactionState
	PostOperationHooks []string
}

func (k *Fsm) GetState(stateName string) (*TransactionState, error) {
	if _, ok := k.States[stateName]; !ok {
		return nil, errors.New(fmt.Sprintf("unable to find state %s", stateName))
	}

	return k.States[stateName], nil
}

func (k *Fsm) AddFSM(name string, handlerFunc StateActionFunc, hookAction HookRequestFunc) {
	newState := new(TransactionState)
	newState.StateName = name
	newState.StateAction = handlerFunc
	newState.HookAction = hookAction
	k.States[name] = newState
}

func (k *Fsm) GoToStateFrom(from string, to string, arg interface{}) (interface{}, error) {
	//get the state
	state, err := k.GetState(from)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid origin state %s", from))
	}

	if _, ok := state.Transitions[to]; !ok {
		return nil, errors.New(fmt.Sprintf("unable to go to %s from %s", to, from))
	}
	//end of validation

	//set the fsm object
	nextTransition := state.Transitions[to]
	operationRes, err := nextTransition.State.StateAction(nextTransition.State, arg)
	//return nextState.HookAction(nextState, operationRes, err)
	if state.HookAction == nil{
		return operationRes, err
	}
	return state.HookAction(nextTransition.PostOperationHooks, operationRes, err)
}

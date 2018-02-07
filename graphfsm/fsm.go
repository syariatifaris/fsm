package graphfsm

import (
	"errors"
	"fmt"
)

const InitialStateName = "-1"

type StateOrderAction interface {
	StateDefaultHandler(fsmObj *TransactionFSMObj, totalArgs int, data interface{})
	EmptyStateHandler(fsmObj *TransactionFSMObj, totalArgs int, data interface{})
}

type StateActionFunc func(*TransactionFSMObj, interface{}) (interface{}, error)

type TransactionFSMState struct {
	StateName          string
	StateAction        StateActionFunc
	PossibleNextStates map[string]*TransactionFSMState
	NextLList          *TransactionFSMState
}

type TransactionFSMObj struct {
	//Ptr to linked list of fsm state structure
	FSMBase *TransactionFSMState
	//Ptr to Current fsm state
	CurrentFSM *TransactionFSMState
	//Name of current fsm
	CurrentFSMStateName string
	//Value to be passed to next
	Arg interface{}
}

//Initialize the fsm
func InitFSM(fsm *TransactionFSMObj) {
	fsm = &TransactionFSMObj{}
}

//Execution of the next state
func nextState(fsm *TransactionFSMObj) error {
	tmp := fsm.FSMBase
	if fsm.FSMBase == nil || fsm.CurrentFSMStateName == "" {
		return errors.New("fsm has not been initialized")
	}

	for (tmp.StateName != fsm.CurrentFSMStateName) && (tmp != nil) {
		tmp = tmp.NextLList
	}

	if tmp == nil {
		return errors.New("unable to find next state")
	}

	tmp.StateAction(fsm, fsm.Arg)
	return nil
}

func ExecuteMain(fsm *TransactionFSMObj) {
	for nextState(fsm) == nil {

	}
}

func AddFSM(fsm *TransactionFSMObj, name string, handlerFunc StateActionFunc) {
	tmpState := fsm.FSMBase
	newState := new(TransactionFSMState)

	for tmpState.NextLList != nil {
		tmpState = tmpState.NextLList
	}

	newState.StateName = name
	newState.StateAction = handlerFunc
	newState.NextLList = nil
	tmpState.NextLList = newState
}

func RemoveFSM(fsm *TransactionFSMObj, name string) error {
	if name != InitialStateName {
		return errors.New("unable to remove initial state")
	}

	var tmpState *TransactionFSMState

	for tmpState.NextLList != nil && tmpState.StateName != name {
		tmpState = tmpState.NextLList
	}

	if tmpState == nil {
		return errors.New(fmt.Sprintf("cannot find state %s", name))
	}

	tmpState.NextLList = nil
	return nil
}

func changeFSMState(fsm *TransactionFSMObj, stateName string, arg interface{}) error {
	tmpState := fsm.FSMBase

	for tmpState != nil && tmpState.StateName != stateName {
		tmpState = tmpState.NextLList
	}

	if tmpState == nil {
		return errors.New(fmt.Sprintf("cannot find state %s", stateName))
	}

	fsm.CurrentFSM = tmpState
	fsm.CurrentFSMStateName = tmpState.StateName
	fsm.Arg = arg
	return nil
}

func GetState(fsm *TransactionFSMObj, stateName string) *TransactionFSMState {
	tmpState := fsm.FSMBase

	for tmpState != nil && tmpState.StateName != stateName {
		tmpState = tmpState.NextLList
	}

	return tmpState
}

func InitialFSM(fsm *TransactionFSMObj, name string, handlerFunc StateActionFunc) {
	fsm.FSMBase = &TransactionFSMState{
		StateName:   name,
		StateAction: handlerFunc,
		NextLList:   nil,
	}
	fsm.CurrentFSM = fsm.FSMBase
	fsm.CurrentFSMStateName = fsm.FSMBase.StateName
}

func goToState(fsm *TransactionFSMObj, name string, arg interface{}) (interface{}, error) {
	err := changeFSMState(fsm, name, arg)
	if err != nil {
		return nil, err
	}

	return fsm.CurrentFSM.StateAction(fsm, fsm.Arg)
}

func GoToStateFrom(fsm *TransactionFSMObj, from string, to string, arg interface{}) (interface{}, error) {
	//go to initial first
	err := changeFSMState(fsm, InitialStateName, arg)
	if err != nil {
		return nil, err
	}

	//validation take place here
	state := GetState(fsm, from)
	if _, ok := state.PossibleNextStates[to]; !ok {
		return nil, errors.New(fmt.Sprintf("unable to go to %s from %s", to, from))
	}
	//end of validation

	//set the fsm object
	nextState := state.PossibleNextStates[to]
	fsm.CurrentFSM = nextState
	fsm.CurrentFSMStateName = nextState.StateName
	fsm.Arg = arg

	return nextState.StateAction(fsm, arg)
}

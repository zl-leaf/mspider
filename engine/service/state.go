package service
import(
    "fmt"
)

const (
    WorkingStateCode = 1
    StopStateCode = 2
    FreeStateCode = 3
)

type IState interface {
    Code() int
    Next(stateCode int) (IState, error)
}

type BaseState struct {
    StateCode int
}
func (this BaseState) Code() int {
    return this.StateCode
}

type WorkingState struct{
    BaseState
}
func (this WorkingState) Next(stateCode int) (nextState IState, err error) {
    err = checkCode(stateCode)
    if err != nil {
        return
    }
    if stateCode != StopStateCode {
        curStateDesc := description(StopStateCode)
        nextStateDesc := description(stateCode)
        err = fmt.Errorf("%s can not change to %s", curStateDesc, nextStateDesc)
        return
    }
    nextState = createState(stateCode)
    return
}

type StopState struct{
    BaseState
}
func (this StopState) Next(stateCode int) (nextState IState, err error) {
    err = checkCode(stateCode)
    if err != nil {
        return
    }
    if stateCode != FreeStateCode {
        curStateDesc := description(StopStateCode)
        nextStateDesc := description(stateCode)
        err = fmt.Errorf("%s can not change to %s", curStateDesc, nextStateDesc)
        return
    }
    nextState = createState(stateCode)
    return
}

type FreeState struct{
    BaseState
}
func (this FreeState) Next(stateCode int) (nextState IState, err error) {
    err = checkCode(stateCode)
    if err != nil {
        return
    }
    if stateCode != WorkingStateCode {
        curStateDesc := description(FreeStateCode)
        nextStateDesc := description(stateCode)
        err = fmt.Errorf("%s can not change to %s", curStateDesc, nextStateDesc)
        return
    }
    nextState = createState(stateCode)
    return
}

func checkCode(stateCode int) error {
    allowCodes := [3]int{WorkingStateCode, StopStateCode, FreeStateCode}
    allow := false
    for _,v := range allowCodes {
        if stateCode == v {
            allow = true
            break
        }
    }
    if !allow {
        return fmt.Errorf("stateCode:%d is undefined", stateCode)
    }
    return nil
}

func description(stateCode int) string {
    var desc string
    switch stateCode {
    case WorkingStateCode:
        desc = "working state"
        break
    case StopStateCode:
        desc = "stop state"
        break
    case FreeStateCode:
        desc = "free state"
        break
    default:
        desc = "stateCode is undefined"
    }
    return desc
}

func createState(stateCode int) IState {
    switch stateCode {
    case WorkingStateCode:
        workingState := WorkingState{}
        workingState.StateCode = WorkingStateCode
        return workingState
    case StopStateCode:
        stopState := StopState{}
        stopState.StateCode = StopStateCode
        return stopState
    case FreeStateCode:
        freeState := FreeState{}
        freeState.StateCode = FreeStateCode
        return freeState
    }
    return nil
}
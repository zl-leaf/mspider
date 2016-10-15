package engine
import(
    "testing"
)

func TestCodeFunc(t *testing.T) {
    workingState := createState(WorkingStateCode)
    stopState := createState(StopStateCode)
    freeState := createState(FreeStateCode)

    if workingState.Code() != WorkingStateCode{
        t.Error("code func error")
    }
    if stopState.Code() != StopStateCode {
        t.Error("code func error")
    }
    if freeState.Code() != FreeStateCode {
        t.Error("code func error")
    }
}

func TestWorkingState(t *testing.T) {
    workingState := createState(WorkingStateCode)
    _, err := workingState.Next(StopStateCode)
    if err != nil {
        t.Error("working next error, it should can change to stop state")
    }

    _, err = workingState.Next(FreeStateCode)
    if err == nil {
        t.Error("working next error, it should not can change to free state")
    }
}

func TestStopState(t *testing.T) {
    stopState := createState(StopStateCode)
    _, err := stopState.Next(FreeStateCode)
    if err != nil {
        t.Error("stop next error, it should can change to free state")
    }

    _, err = stopState.Next(WorkingStateCode)
    if err == nil {
        t.Error("working next error, it should not can change to working state")
    }
}

func TestFreeState(t *testing.T) {
    freeState := createState(FreeStateCode)
    _, err := freeState.Next(WorkingStateCode)
    if err != nil {
        t.Error("working next error, it should can change to working state")
    }

    _, err = freeState.Next(StopStateCode)
    if err == nil {
        t.Error("working next error, it should not can change to stop state")
    }
}
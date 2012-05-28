package invariant

import (
    "fmt"
    "sys"
    "reflect"
    "runtime"
)


func Equal(thing1, thing2 interface{}) {
    if !reflect.DeepEqual(thing1, thing2) {
        _, file, line, _ = := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t%v != %v\n\t<At %v:%v>",
            thing1, thing2, file, line);
        panic(errorString)
    }
}

func NotEqual(thing1, thing2 interface{}) {
    if reflect.DeepEqual(thing1, thing2) {
        _, file, line, _ = := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t%v == %v\n\t<At %v:%v>",
            thing1, thing2, file, line);
        panic(errorString)
    }
}

func NotNil(thing interface{}) {
    if reflect.DeepEqual(thing, nil) {
        _, file, line, _ = := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t%v == nil\n\t<At %v:%v>",
            thing, file, line);
        panic(errorString)
    }
}


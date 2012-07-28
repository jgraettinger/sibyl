package invariant

import (
    "fmt"
    "reflect"
    "runtime"
)


func Equal(thing1, thing2 interface{}, args ...interface{}) {
    if !reflect.DeepEqual(thing1, thing2) {
        _, file, line, _ := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t%v != %v\n\t<At %v:%v>",
            thing1, thing2, file, line)
        panic(errorString)
    }
}

func NotEqual(thing1, thing2 interface{}, args ...interface{}) {
    if reflect.DeepEqual(thing1, thing2) {
        _, file, line, _ := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t%v == %v\n\t<At %v:%v>",
            thing1, thing2, file, line)
        panic(errorString)
    }
}

func NotNil(thing interface{}, args ...interface{}) {
    if reflect.ValueOf(thing).IsNil() {
        _, file, line, _ := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t%v == nil\n\t<At %v:%v>",
            thing, file, line)
        panic(errorString)
    }
}

func IsNil(thing interface{}, args ...interface{}) {
    if !reflect.ValueOf(thing).IsNil() {
        _, file, line, _ := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t%v != nil\n\t<At %v:%v>",
            thing, file, line)
        panic(errorString)
    }
}

func IsTrue(result bool, args ...interface{}) {
    if !result {
        _, file, line, _ := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t<At %v:%v>", file, line)
        panic(errorString)
    }
}

func IsFalse(result bool, args ...interface{}) {
    if result {
        _, file, line, _ := runtime.Caller(2)

        errorString := fmt.Sprintf(
            "Invariant Violation:\n\t<At %v:%v>", file, line)
        panic(errorString)
    }
}

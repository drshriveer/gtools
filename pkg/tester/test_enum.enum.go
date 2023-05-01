package tester

import (
    "fmt"
)

func (e MyEnum) IsValid(in MyEnum) bool {
    return true
}

func (e MyEnum) Values() []MyEnum {
    return []MyEnum{
        IntentionallyNegative,
        EnumOneComplicationZero,
        EnumThreeComplicationOne,
        EnumTwoComplicationOne,
        EnumTwoComplicationZero,
        UNSET,
        EnumOneComplicationOne,
        EnumThreeComplicationThree,
        EnumTwoComplicationThree,
        EnumTwoComplicationTwo,
        ValueOne,
        EnumOneComplicationTwo,
        ValueTwo,
        EnumThreeComplicationTwo,
        EnumThreeComplicationZero,
        ValueSeven,
    }
}

func (e MyEnum) String() string {
    switch e {
    case IntentionallyNegative:
        return "IntentionallyNegative"
    case EnumOneComplicationZero:
        return "EnumOneComplicationZero"
    case EnumOneComplicationOne:
        return "EnumOneComplicationOne"
    case EnumOneComplicationTwo:
        return "EnumOneComplicationTwo"
    case EnumThreeComplicationTwo:
        return "EnumThreeComplicationTwo"
    default:
        return fmt.Sprintf("UndefinedMyEnum:%d", e)
    }
}

func (e MyEnum) ParseString(text string) (MyEnum, bool) {
    switch text {
    case "IntentionallyNegative":
        return IntentionallyNegative, true
    case "EnumOneComplicationZero":
        return EnumOneComplicationZero, true
    case "EnumThreeComplicationOne":
        return EnumThreeComplicationOne, true
    case "EnumTwoComplicationOne":
        return EnumTwoComplicationOne, true
    case "EnumTwoComplicationZero":
        return EnumTwoComplicationZero, true
    case "UNSET":
        return UNSET, true
    case "EnumOneComplicationOne":
        return EnumOneComplicationOne, true
    case "EnumThreeComplicationThree":
        return EnumThreeComplicationThree, true
    case "EnumTwoComplicationThree":
        return EnumTwoComplicationThree, true
    case "EnumTwoComplicationTwo":
        return EnumTwoComplicationTwo, true
    case "ValueOne":
        return ValueOne, true
    case "EnumOneComplicationTwo":
        return EnumOneComplicationTwo, true
    case "ValueTwo":
        return ValueTwo, true
    case "EnumThreeComplicationTwo":
        return EnumThreeComplicationTwo, true
    case "EnumThreeComplicationZero":
        return EnumThreeComplicationZero, true
    case "ValueSeven":
        return ValueSeven, true
    default:
        return 0, false
    }
}
import Foundation

indirect enum Value {
    
    case unit
    case bool
    case int
    case float
    case string
    case list(Value)
    case tuple([Value])
    case reference(ValueHolder)
    case function(Function)
    case primitive(Primitive)
    case exception(Exception)
    
}

class ValueHolder {
    
    var value: Value
    
    init(value: Value) {
        self.value = value
    }
    
}

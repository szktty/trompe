import Foundation

indirect enum Type {

    case unit
    case bool
    case int
    case float
    case string
    case list(Type)
    case tuple([Type])
    case reference(Type)
    case function(arguments: [Type], `return`: Type)
    case exception(ExceptionType)
    
}

class ExceptionType {
    
}

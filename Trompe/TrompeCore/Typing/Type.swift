import Foundation

// type annotation
class TypeAnnot {
    
    var location: Location?
    var type: Type
    
    init(location: Location? = nil, type: Type) {
        self.location = location
        self.type = type
    }
    
    static var bool: TypeAnnot = TypeAnnot(type: Type.bool)
    static var int: TypeAnnot = TypeAnnot(type: Type.int)
    
    static func alias(location: Location? = nil,
                      namePath: NamePath,
                      annot: TypeAnnot) -> TypeAnnot {
        return TypeAnnot(location: location,
                         type: .alias(namePath: namePath, annot: annot))
    }
    
}

indirect enum Type: Equatable {

    case `var`(TypeVar)
    case app(const: TypeConst, args: [TypeAnnot])
    case poly(vars: [TypeVar], annot: TypeAnnot)
    case meta(MetaType)
    case alias(namePath: NamePath, annot: TypeAnnot)
 
    static var bool: Type = Type.app(const: .bool, args: [])
    static var int: Type = Type.app(const: .int, args: [])
    
    static func == (lhs: Type, rhs: Type) -> Bool {
        // TODO
        return true
    }
    
}

class TypeVar {
    
}

enum TypeConst: Equatable {

    case unit
    case bool
    case int
    case float
    case string
    case list(Type)
    case tuple([Type])
    case option(Type)
    case ref(Type)
    case fun(args: [Type], `return`: Type)
    case exn(ExnType)
    
    static func == (lhs: TypeConst, rhs: TypeConst) -> Bool {
        switch (lhs, rhs) {
        case (.unit, .unit),
             (.bool, .bool),
             (.int, .int),
             (.float, .float),
             (.string, .string):
            return true
        case (.list(let a), .list(let b)),
             (.option(let a), .option(let b)),
             (.ref(let a), .ref(let b)):
            return a == b
        case (.tuple(let a), .tuple(let b)):
            return a == b
        default:
            // TODO
            return false
        }
    }
    
}

class MetaType {
    
    var annot: TypeAnnot?
    
    var isUndefined: Bool {
        get { return annot == nil }
    }
    
}

// exception type
class ExnType {
    
    var namePath: NamePath
    
    init(namePath: NamePath) {
        self.namePath = namePath
    }
    
}

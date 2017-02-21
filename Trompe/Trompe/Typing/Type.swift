import Foundation

// type annotation
struct TypeAnnot {
    
    var location: Location?
    var type: Type
    
    init(location: Location? = nil, type: Type) {
        self.location = location
        self.type = type
    }
    
    static var int: TypeAnnot = TypeAnnot(type: Type.int)
    
    static func alias(location: Location? = nil,
                      namePath: NamePath,
                      annot: TypeAnnot) -> TypeAnnot {
        return TypeAnnot(location: location,
                         type: .alias(namePath: namePath, annot: annot))
    }
    
}

indirect enum Type: Comparable {

    case `var`(TypeVar)
    case app(const: TypeConst, args: [TypeAnnot])
    case poly(vars: [TypeVar], annot: TypeAnnot)
    case meta(MetaType)
    case alias(namePath: NamePath, annot: TypeAnnot)
 
    static var int: Type = Type.app(const: .int, args: [])
    
    static func < (lhs: Type, rhs: Type) -> Bool {
        // TODO
        return true
    }
    
    static func == (lhs: Type, rhs: Type) -> Bool {
        // TODO
        return true
    }
    
}

class TypeVar {
    
}

enum TypeConst {

    case unit
    case bool
    case int
    case float
    case string
    case list(Type)
    case tuple([Type])
    case ref(Type)
    case fun(args: [Type], `return`: Type)
    case exn(ExnType)
    
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

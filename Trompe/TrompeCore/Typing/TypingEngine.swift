import Foundation

enum TypeError: Error {
    case mismatch(expected: TypeAnnot, actual: TypeAnnot)
}

class TypingEngine {
    
    var moduleManager: TypeModuleManager = TypeModuleManager()

    func unify(expected: TypeAnnot, actual: TypeAnnot) throws {
        if expected.type == actual.type {
            return
        }
        
        let mismatch = TypeError.mismatch(expected: expected, actual: actual)
        switch (expected.type, actual.type) {
        case (.meta(let meta1), .meta(let meta2)):
            switch (meta1.annot, meta2.annot) {
            case (nil, let annot):
                meta1.annot = annot
            case (let annot, nil):
                meta2.annot = annot
            default:
                try unify(expected: meta1.annot!, actual: meta2.annot!)
            }
    
        case (.app(const: let const1, args: let args1),
              .app(const: let const2, args: let args2)):
            if const1 != const2 || args1.count != args2.count {
                throw mismatch
            }
            for i in args1.indices {
                let arg1 = args1[i]
                let arg2 = args2[i]
                try unify(expected: arg1, actual: arg2)
            }
            
        default:
            // TODO
            return
        }
    }
    
}

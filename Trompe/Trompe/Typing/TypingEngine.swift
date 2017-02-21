import Foundation

enum TypeError: Error {
    case mismatch(expected: TypeAnnot, actual: TypeAnnot)
}

class TypingEngine {
    
    func unify(expected: TypeAnnot, actual: TypeAnnot) throws {
        if expected.type == actual.type {
            return
        }
        
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
        default:
            // TODO
            return
        }
    }
    
}

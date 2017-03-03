import Foundation

class Env<T> {
    
    weak var parent: Env<T>?
    var attributes: [String: T] = [:]
    var imports: [Env<T>] = []
    
    required init(parent: Env<T>? = nil) {
        self.parent = parent
    }
    
}

typealias ValueEnv = Env<Value>
typealias TypeEnv = Env<Type>

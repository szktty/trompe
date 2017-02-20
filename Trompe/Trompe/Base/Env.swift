import Foundation

class Env<T> {
    
    weak var parentEnv: Env<T>?
    var attributes: [String: T] = [:]
    var importedEnvs: [Env<T>] = []
    
    init(parentEnv: Env<T>? = nil) {
        self.parentEnv = parentEnv
    }
    
}

typealias ValueEnv = Env<Value>
typealias TypeEnv = Env<Type>

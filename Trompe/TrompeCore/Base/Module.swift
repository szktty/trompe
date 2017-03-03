import Foundation

class Module<T> {
    
    weak var parent: Module<T>?
    var name: String?
    var env: Env<T>
    var submodules: [Module<T>] = []
    var imports: [Module<T>] = []
    
    init(parent: Module<T>? = nil, name: String? = nil) {
        self.parent = parent
        self.name = name
        if let parent = parent {
            env = Env(parent: parent.env)
        } else {
            env = Env()
        }
    }
    
    func find(name: String) -> Module<T>? {
        if let found = (submodules.first { m in return m.name == name }) {
            return found
        }
        for import_ in imports {
            if let found = import_.find(name: name) {
                return found
            }
        }
        return nil
    }
    
    func find(namePath: NamePath) -> Module<T>? {
        assertionFailure()
        return nil
    }
    
}

typealias ValueModule = Module<Value>
typealias TypeModule = Module<Type>

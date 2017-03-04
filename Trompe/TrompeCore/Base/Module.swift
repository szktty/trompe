import Foundation

class Module<T>: Equatable {
    
    static func == (lhs: Module<T>, rhs: Module<T>) -> Bool {
        return ObjectIdentifier(lhs) == ObjectIdentifier(rhs)
    }
    
    weak var parent: Module<T>?
    var name: String?
    var env: Env<T> = Env()
    var submodules: [Module<T>] = []
    var imports: [Module<T>] = []
    
    var isRoot: Bool {
        get { return parent == nil }
    }
    
    init(name: String? = nil) {
        self.name = name
    }
    
    func add(module: Module<T>) {
        module.parent = self
        module.env.parent = env
        submodules.append(module)
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

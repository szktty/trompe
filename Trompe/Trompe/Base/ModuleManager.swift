import Foundation

class ModuleManager<T> {

    var root: Module<T> = Module()
    
}

typealias ValueModuleManager = ModuleManager<Value>
typealias TypeModuleManager = ModuleManager<Type>

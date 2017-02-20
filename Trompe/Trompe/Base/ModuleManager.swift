import Foundation

var valueModuleManager: ModuleManager<Value> = ModuleManager()
var typeModuleManager: ModuleManager<Type> = ModuleManager()

class ModuleManager<T> {

    var rootModule: Module<T> = Module()
    
}

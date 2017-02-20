import Foundation

class Module<T> {
    
    weak var parentModule: Module<T>?
    var env: Env<T>
    var submodules: [Module<T>] = []
    var importedModules: [Module<T>] = []
    
    init(parentModule: Module<T>? = nil) {
        self.parentModule = parentModule
        if let parent = parentModule {
            env = Env(parentEnv: parent.env)
        } else {
            env = Env()
        }
    }
    
}

import Foundation

class Interpreter {
    
    var moduleManager: ValueModuleManager = ValueModuleManager()

    func evaluate(context: Context, env: ValueEnv, node: Node) -> (ValueEnv, Value) {
        assertionFailure()
        return (ValueEnv(), .unit)
    }
    
}

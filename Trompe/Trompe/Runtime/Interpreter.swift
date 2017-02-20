import Foundation

class Interpreter {
    
    func evaluate(context: Context, env: ValueEnv, node: Node) -> (ValueEnv, Value) {
        assertionFailure()
        return (ValueEnv(), .unit)
    }
    
}

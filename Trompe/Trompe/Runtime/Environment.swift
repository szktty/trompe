import Foundation

class Environment<T> {
    
    var parent: Environment?
    var attributes: [String: T] = [:]
    var importedEnvs: [Environment] = []
    
}

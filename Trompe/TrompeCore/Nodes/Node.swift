import Foundation

protocol Node {
    
    // var location: Location { get }
    
}

indirect enum SeparatedList<Element, Separator> {
    
    case Cons(element: Element, next: SeparatedList<Element, Separator>)
    case Nil
    
}

class Token {
    
    var location: Location?
    var contents: String?
    
}

typealias NodeList = SeparatedList<Node, Token>

class UnaryNode: Node {

    var token: Token
    
    init(token: Token) {
        self.token = token
    }
    
}

class EnclosedNode<Contents>: Node {
    
    var open: Token?
    var contents: Contents?
    var close: Token?
    
}

class FunctionNode: Node {
    
    var prefix: Node?
    var arguments: EnclosedNode<NodeList>?
    
}

class ReferencePathNode: Node {
    
    var path: NodeList?
    
}

class BooleanNode: UnaryNode {}
class StringNode: UnaryNode {}
class IntegerNode: UnaryNode {}
class FloatNode: UnaryNode {}
class TupleNode: EnclosedNode<NodeList> {}
class ListNode: EnclosedNode<NodeList> {}

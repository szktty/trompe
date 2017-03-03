import Foundation

class Position {
    
    static var zero: Position = Position(line: 0, column: 0, offset: 0)
    
    var line: Int
    var column: Int
    var offset: Int

    init(line: Int, column: Int, offset: Int) {
        self.line = line
        self.column = column
        self.offset = offset
    }
    
}

class Location {
    
    var start: Position
    var end: Position
    
    var length: Int {
        get {
            // TODO
            return 0
        }
    }
    
    init(start: Position, end: Position) {
        self.start = start
        self.end = end
    }
    
}

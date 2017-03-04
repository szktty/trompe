import Foundation

struct Position: Equatable {
    
    static var zero: Position = Position(line: 0, column: 0, offset: 0)
    
    var line: Int
    var column: Int
    var offset: Int

    init(line: Int, column: Int, offset: Int) {
        self.line = line
        self.column = column
        self.offset = offset
    }
    
    static func == (lhs: Position, rhs: Position) -> Bool {
        return lhs.line == rhs.line &&
            lhs.column == rhs.column &&
            lhs.offset == rhs.offset
    }
    
}

class Location: Equatable {
    
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
    
    static func == (lhs: Location, rhs: Location) -> Bool {
        return lhs.start == rhs.start && lhs.end == rhs.end
    }
    
}

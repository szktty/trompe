import XCTest
@testable import TrompeCore

class TypingTests: XCTestCase {

    var engine: TypingEngine!

    override func setUp() {
        super.setUp()
        // Put setup code here. This method is called before the invocation of each test method in the class.
        engine = TypingEngine()
    }
    
    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
        super.tearDown()
    }

    func testBasicOK() {
        XCTAssertNotNil(try? engine.unify(expected: .int, actual: .int))
    }
    
    func testBasicFail() {
        XCTAssertNil(try? engine.unify(expected: .int, actual: .bool))
    }

    func testPerformanceExample() {
        // This is an example of a performance test case.
        self.measure {
            // Put the code you want to measure the time of here.
        }
    }

}

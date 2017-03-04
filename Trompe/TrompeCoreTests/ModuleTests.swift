import XCTest
@testable import TrompeCore

typealias TestModuleManager = ModuleManager<String>
typealias TestModule = Module<String>

class ModuleTests: XCTestCase {

    var manager: TestModuleManager!
    var m_a: TestModule!
    var m_b: TestModule!
    var m_c: TestModule!
    
    override func setUp() {
        super.setUp()
        // Put setup code here. This method is called before the invocation of each test method in the class.
        manager = TestModuleManager()
        m_a = TestModule(name: "A")
        m_b = TestModule(name: "B")
        m_c = TestModule(name: "C")
        m_a.add(module: m_b)
        m_b.add(module: m_c)
        manager.root = m_a
    }
    
    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
        super.tearDown()
    }

    func testParent() {
        XCTAssert(m_a.parent == nil)
        XCTAssert(m_b.parent == m_a!)
        XCTAssert(m_c.parent == m_b!)
    }

    func testFind() {
        XCTAssert(m_a.find(name: "B") == m_b)
        XCTAssert(m_b.find(name: "C") == m_c)
        XCTAssert(m_a.find(name: "Z") == nil)
    }
    
    func testPerformanceExample() {
        // This is an example of a performance test case.
        self.measure {
            // Put the code you want to measure the time of here.
        }
    }

}

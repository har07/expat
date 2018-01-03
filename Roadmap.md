## Create XML Processing Library For Go Based On Python's ElementTree ##

# Create Thin Wrapper For C's Expat Library #
- C wrapper that makes calling Expat functions from *cgo* easier
- Go wrapper that makes calling Expat functions from *Go* easier (no need for import "C" at all)

# Create Go XMLParser Struct/Library #
- TODO (Mimic Python's XMLParser class):
    - `Feed(data string)`: receive chunk of XML string. Call Expat XML_Parse with finish=false
    - `Close()`: Call `XML_Parse("", true)` and return error if any. [N] Close TreeBuilder.
    - `Default(text string)`: default handler for any unhandled chunk of XML
    - `Start(tag string, attr map[string]string)`: [N] Call start element handler of TreeBuilder.
    - `End(tag string)`: [N] Call end element handler of TreeBuilder.
    - `RaiseError(??)`
    - `FixName(key string)`
    - `SetEvents(??)`
    - `Create(encoding string, target TreeBuilder)`
    - *[N]: indicates future feature

# Create Go TreeBuilder Struct/Library #
- This is a builder for class Element [v]
    - `Init()`
    - `Close()`
    - `flush()`
    - `Data()`
    - `Start()`
    - `End()`

# Create Go ElementTree Library #
# Create Go Element Class #
- For V0.X, create a very simple Element class implementation, just to enable us to put everything together and test (parse+tree builder+element class) [v]


# The Roadmap #
V0.
- minimal codes to the steps above
- skip: find, findall
- implement: GetElementsByTagName, core functions for event-based API

V0.
- improve architecture if possible
- add more event-based API if any left

V0.
- find, findall: maybe based on Go etree library

V1.
- review implementation for performance, memory efficiency, encoding handling

V2.
- implement better XPath support
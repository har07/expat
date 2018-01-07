## Create XML Processing Library For Go Based On Python's ElementTree ##

# Create Thin Wrapper For C's Expat Library #
- C wrapper that makes calling Expat functions from *cgo* easier
- Go wrapper that makes calling Expat functions from *Go* easier (no need for import "C" at all)

# Create Go XMLParser Struct/Library #

# Create Go TreeBuilder Struct/Library #

# Create Go ElementTree Library #
# Create Go Element Class #
- For V0.X, create a very simple Element class implementation, just to enable us to put everything together and test (parse+tree builder+element class) [v]

# Current #
- Collect namespaces info on ElementTree creation (on parsing)
- Figure out why C.free() threw exception

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
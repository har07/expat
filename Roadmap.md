## Create XML Processing Library For Go Based On Python's ElementTree ##

# Create Thin Wrapper For C's Expat Library #
- C wrapper that makes calling Expat functions from *cgo* easier
- Go wrapper that makes calling Expat functions from *Go* easier (no need for import "C" at all)
    - read attrib [v]

# Create Go XMLParser Struct/Library #

# Create Go ElementTree Library #


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
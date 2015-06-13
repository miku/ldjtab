README
======

ldjtab extract a list of values from a given key and the line number of the
document from a line-delimited JSON file.

Install with:

    $ go get github.com/miku/ldjtab/cmd/...

Example:

    $ cat test.txt
    {"id": 1, "name": "A"}
    {"id": 2, "name": "B"}
    {"id": 3, "name": "C"}

    $ ldjtab -key name test.txt
    A    0000000001
    B    0000000002
    C    0000000003

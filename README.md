README
======

ldjtab extract a list of values from a given key and the line number of the
document from a line-delimited JSON file. It's reasonably fast and will try to
utilize all cores.

For a more interesting extractions, use [jq](http://stedolan.github.io/jq/).

Install with:

    $ go get github.com/miku/ldjtab/cmd/...

Note: There are two binaries at the moment: `ldjtab` and `uldjtab`. `ldjtab` works
in parallel, is a bit faster, but does not guarantee order. `uldjtab` is a bit
slower, lighter on resources and preserves order.

Example:

    $ cat fixtures/test.txt
    {"id": 1, "name": "A"}
    {"id": 2, "name": "B"}
    {"id": 3, "name": "C"}

    $ ldjtab -key name fixtures/test.txt
    A    0000000001
    B    0000000002
    C    0000000003

    $ uldjtab -key name fixtures/test.txt
    A    0000000001
    B    0000000002
    C    0000000003

Use case
--------

JSON file compaction.

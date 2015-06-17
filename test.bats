#!/usr/bin/env bats

LDJTAB="go run cmd/ldjtab/main.go"
TESTFILE_CROSSREF=fixtures/crossref.10.ldj

@test "extract URL, line numbers" {
  result="$($LDJTAB -key URL $TESTFILE_CROSSREF | wc -l)"
  [ "$result" -eq 10 ]
}

@test "extract URL, values" {
  result="$(go run cmd/ldjtab/main.go -key URL fixtures/crossref.10.ldj | head -1| cut -f1 | tr -d '\n')"
  [ "$result" = "http://dx.doi.org/10.1590/s0100-41582006000200010" ]

  result="$(go run cmd/ldjtab/main.go -key URL fixtures/crossref.10.ldj | head -1 | cut -f2| tr -d '\n')"
  [ "$result" = "0000000001" ]
}

@test "extract URL, padding" {
  result="$(go run cmd/ldjtab/main.go -padlength 0 -key URL fixtures/crossref.10.ldj | head -1| cut -f1 | tr -d '\n')"
  [ "$result" = "http://dx.doi.org/10.1590/s0100-41582006000200010" ]

  result="$(go run cmd/ldjtab/main.go -padlength 0 -key URL fixtures/crossref.10.ldj | head -1 | cut -f2| tr -d '\n')"
  [ "$result" = "1" ]
}

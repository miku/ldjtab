#!/usr/bin/env bats

@test "ldjtab: extract URL, line numbers" {
  result="$(go run cmd/ldjtab/main.go -key URL fixtures/crossref.10.ldj | wc -l)"
  [ "$result" -eq 10 ]
}

@test "ldjtab: extract URL, values" {
  result="$(go run cmd/ldjtab/main.go -key URL fixtures/crossref.10.ldj | head -1| cut -f1 | tr -d '\n')"
  [ "$result" = "http://dx.doi.org/10.1590/s0100-41582006000200010" ]

  result="$(go run cmd/ldjtab/main.go -key URL fixtures/crossref.10.ldj | head -1 | cut -f2| tr -d '\n')"
  [ "$result" = "0000000001" ]
}

@test "ldjtab: extract URL, padding" {
  result="$(go run cmd/ldjtab/main.go -padlength 0 -key URL fixtures/crossref.10.ldj | head -1| cut -f1 | tr -d '\n')"
  [ "$result" = "http://dx.doi.org/10.1590/s0100-41582006000200010" ]

  result="$(go run cmd/ldjtab/main.go -padlength 0 -key URL fixtures/crossref.10.ldj | head -1 | cut -f2| tr -d '\n')"
  [ "$result" = "1" ]
}

@test "uldjtab: extract URL, line numbers" {
  result="$(go run cmd/uldjtab/main.go -key URL fixtures/crossref.10.ldj | wc -l)"
  [ "$result" -eq 10 ]
}

@test "uldjtab: extract URL, values" {
  result="$(go run cmd/uldjtab/main.go -key URL fixtures/crossref.10.ldj | head -1| cut -f1 | tr -d '\n')"
  [ "$result" = "http://dx.doi.org/10.1590/s0100-41582006000200010" ]

  result="$(go run cmd/uldjtab/main.go -key URL fixtures/crossref.10.ldj | head -1 | cut -f2| tr -d '\n')"
  [ "$result" = "0000000001" ]
}

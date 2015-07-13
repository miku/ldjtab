SHELL = /bin/bash
TARGETS = ldjtab uldjtab

all: uldjtab ldjtab

uldjtab: deps imports
	go build -o uldjtab cmd/uldjtab/main.go

ldjtab: deps imports
	go build -o ldjtab cmd/ldjtab/main.go

deps:
	go get ./...

imports:
	goimports -w .

clean:
	rm -f $(TARGETS)
	rm -f ldjtab_*deb
	rm -f ldjtab-*rpm
	rm -rf ./packaging/deb/ldjtab/usr

test:
	bats test.bats

deb: $(TARGETS)
	mkdir -p packaging/deb/ldjtab/usr/sbin
	cp $(TARGETS) packaging/deb/ldjtab/usr/sbin
	cd packaging/deb && fakeroot dpkg-deb --build ldjtab .
	mv packaging/deb/ldjtab_*.deb .

rpm: $(TARGETS)
	mkdir -p $(HOME)/rpmbuild/{BUILD,SOURCES,SPECS,RPMS}
	cp ./packaging/rpm/ldjtab.spec $(HOME)/rpmbuild/SPECS
	cp $(TARGETS) $(HOME)/rpmbuild/BUILD
	./packaging/rpm/buildrpm.sh ldjtab
	cp $(HOME)/rpmbuild/RPMS/x86_64/ldjtab*.rpm .

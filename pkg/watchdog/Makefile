SHELL = /bin/bash

.PHONY: test

# these tests run in isolation by calling go test -run=... or the equivalent.
ISOLATED_TESTS +=
ifdef CI
	ISOLATED_TESTS = TestControl_Isolated \
					 TestIsolatedControlCGroup \
					 TestIsolatedSystemDriven \
					 TestIsolatedHeapDriven
else
	ISOLATED_TESTS = TestControl_Isolated \
					 TestIsolatedControlCGroup \
					 TestIsolatedSystemDriven \
					 TestIsolatedHeapDriven
endif

test: test-binary test-docker

test-binary:
	go test -v ./... # run all the non-isolated tests.
	# foreach does not actually execute each iteration; it expands the text, and it's executed all at once
    # that's why we use && true, to shorcircuit if a test fails.
	$(foreach name,$(ISOLATED_TESTS),TEST_ISOLATED=1 go test -v -test.run=$(name) ./... && ) true

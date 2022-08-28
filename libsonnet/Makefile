gen-lib:
	@echo Generating library
	@rm -rfd 1.0/_gen/*
	@rm 1.0/*.libsonnet
	@python scripts/jsonnet-lib-gen.py --output-dir 1.0/ ../docs/gen/policies/gen.yaml
	@tk fmt 1.0/
	@git add 1.0/*

.PHONY: gen-lib

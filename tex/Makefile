.ONESHELL:

VERSION?=$(shell grep 'Version' ../cli/cmd/root.go | awk '{ print $$4 }' | head -n 1 | tr -d '"')

cls:
	cd latex/csassignments
	latex.exe csassignments.ins

docs:
	cd latex/csassignments
	pdflatex.exe -interaction=nonstopmode csassignments.dtx
	latexmk.exe -c csassignments.dtx
	rm *.glo
	rm *.hd

extract-docs: docs
	cd latex/csassignments
	mv csassignments.pdf ../csassignments-v$(VERSION).pdf

zip: cls docs
	cd latex
	zip -r csassignments-v$(VERSION).zip csassignments

release: zip extract-docs

.PHONY: zip
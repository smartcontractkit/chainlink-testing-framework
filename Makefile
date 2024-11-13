.PHONY: modgraph
modgraph:
	go install github.com/jmank88/gomods@v0.1.5
	go install github.com/jmank88/modgraph@v0.1.0
	./modgraph > go.md

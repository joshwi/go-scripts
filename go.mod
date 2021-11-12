module github.com/joshwi/go-scripts

go 1.16

replace github.com/joshwi/go-scripts/collector => ./collector

replace github.com/joshwi/go-scripts/git => ./git

require (
	github.com/joshwi/go-git v1.0.0
	github.com/joshwi/go-plugins v1.0.0
	github.com/joshwi/go-utils v1.0.3
	github.com/neo4j/neo4j-go-driver v1.8.3
)

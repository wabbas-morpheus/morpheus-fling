module github.com/wabbas-morpheus/morpheus-fling/elasticIng

go 1.20

require (
	github.com/elastic/go-elasticsearch/v7 v7.17.7
	github.com/mitchellh/mapstructure v1.5.0
)

require github.com/wabbas-morpheus/morpheus-fling/rbParse v0.0.0-20240202132330-88e861442fe3

replace github.com/wabbas-morpheus/morpheus-fling/rbParse => ../rbParse

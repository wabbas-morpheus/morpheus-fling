module morpheus-fling

go 1.20

replace github.com/wabbas-morpheus/morpheus-fling/fileReader => ./fileReader

replace github.com/wabbas-morpheus/morpheus-fling/rabbitIng => ./rabbitIng

require (
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/wabbas-morpheus/morpheus-fling/elasticIng v0.0.0-20230317162923-208ce5212288
	github.com/wabbas-morpheus/morpheus-fling/encryptText v0.0.0-20230317162923-208ce5212288
	github.com/wabbas-morpheus/morpheus-fling/fileReader v0.0.0-00010101000000-000000000000
	github.com/wabbas-morpheus/morpheus-fling/portScanner v0.0.0-20230317162923-208ce5212288
	github.com/wabbas-morpheus/morpheus-fling/rabbitIng v0.0.0-00010101000000-000000000000
	github.com/wabbas-morpheus/morpheus-fling/secParse v0.0.0-20230317162923-208ce5212288
	github.com/wabbas-morpheus/morpheus-fling/sysGatherer v0.0.0-20230317162923-208ce5212288
	github.com/zcalusic/sysinfo v0.9.5
)

require (
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/elastic/go-elasticsearch/v7 v7.17.7 // indirect
	github.com/frankban/quicktest v1.14.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/sync v0.1.0 // indirect
)

replace github.com/wabbas-morpheus/morpheus-fling/elasticIng => ./elasticIng

replace github.com/wabbas-morpheus/morpheus-fling/encryptText => ./encryptText

replace github.com/wabbas-morpheus/morpheus-fling/portScanner => ./portScanner

replace github.com/wabbas-morpheus/morpheus-fling/secParse => ./secParse

replace github.com/wabbas-morpheus/morpheus-fling/sysGatherer => ./sysGatherer

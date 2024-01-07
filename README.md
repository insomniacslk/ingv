# ingv

Go library and CLI to fetch earthquake information from the Italian
[INGV](https://ingv.it).

The INGV Earthquake Event API is documented at
https://webservices.ingv.it/swagger-ui/dist/?url=https://ingv.github.io/openapi/fdsnws/event/0.0.1/event.yaml
, which is the API equivalent of visiting the web page http://terremoti.ingv.it/ .


## Command line interface

Located at [cmd/quakes]. Just run `go build`, then `./quakes --help`.

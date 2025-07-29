### 🚀 Running the App

#### Compile and run locally (CLI mode)
- Requirements: go 1.24+ 
- How to run: clone the project then `go build ./cmd/essayanalyzer`, then `./essayanalyzer`
- For help: `./essayanalyzer -h`

#### Use the docker image (use default parameters)

- Requirements: docker
- How to run:  `docker run ghcr.io/alexensimag/firefly:latest`

### Considerations

#### Configurable parameters
- minimum word size (default: `3`)
- number of top words returned (default: `10`)
- number of workers (default: `20`)
- number of tasks loaded in the "task buffer": (default: `200`)


#### Normalization algorithm
Additionally to the requirements, when evaluated words are normalized as such:
- words are lowercased
- punctuations at the beginning and end of words are removed (besides `'` at end of word)
- examples:
    - `the` and `The` are both normalized as `the`
    - `end` and `end.` are both normalized as `end.`
    - `its` and `it's` are respectively normalized as `its` and `it's`
    - `(and` is normalized as `and`
    - `'and'` is normalized as `and'`


### GitHub actions
- "Go PR Check": runs on commit (goimports, go vet, go test)
- "Push Docker image to GHCR": builds the image and pushes it to GHCR

## Possible extensions
- Set the config with environment variables to support parameters in the Docker image
- The normalization algorithm is not perfect (single quoted words are not handled properly) and was built to support English, it should be overridable, in order to support different logics and languages
- The app supports only Engadget essays, it was designed to be extended with an interface `Task` to support different sources and resource types
- Need for metrics to be exposed, with custom metrics showing progression
- Need for the app the be "interruptable" and restart from where it stopped


---



### Run Locally

```bash
go run ./cmd/app
```


![Coverage](https://codecov.io/gh/alexEnsimag/firefly/branch/main/graph/badge.svg)
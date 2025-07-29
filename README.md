![Coverage](https://img.shields.io/badge/coverage-70%25-brightgreen)

## EssayAnalyzer

A simple and efficient tool to cound predefined words in essays with configurable parameters and Docker support.

### 🚀 How to run

#### Compile and Run Locally (CLI Mode)

Requirements: Go 1.24 or higher

1. Clone the repository:
   ```bash
   git clone https://github.com/alexEnsimag/firefly.git
   cd firefly
   ```
2. Build the application:
   ```bash
   go build ./cmd/essayanalyzer
   ```
3. Run the binary:
   ```bash
   ./essayanalyzer
   ./essayanalyzer -h # for help prompt
   ```

#### Run Using Docker (Default Parameters)

Requirements: Docker installed and running

Run:
```bash
docker run ghcr.io/alexensimag/firefly:latest
```

### Configurations

You can configure the following parameters (defaults shown):

- minimum word size: `3`
- number of top words returned: `10`
- number of workers: `20`
- task buffer size: `200`

### Considerations

- The application is cancellable with `Ctrl+C`. Note that ongoing work will be lost on interruption.

- The application uses a custom HTTP client that handles retries and rate limit
    - rate limit is set too `100` requests per second, no burst allowed
    - automatically retries on network errors and HTTP 429 and 5XX 
    - max 5 retries, with exponentiall backoff between 1-5 seconds

- **Normalization algorithm** details:
  - Words are converted to lowercase.
  - Leading and trailing punctuation are removed, except an apostrophe `'` at the end of a word - *needs clearer definition about quotes and punctuation*.
  - Examples:
    - `The` and `the` → `the`
    - `end.` and `end` → `end`
    - `its` and `it's` → `its` and `it's`
    - `(and` → `and`
    - `orders` → `orders'`
    - `'and'` -> `and'`

---

### 🛠️ GitHub Actions

- **Go PR Check:** Runs on every commit, includes `goimports`, `go vet`, and tests.
- **Push Docker Image:** Builds and pushes the Docker image to GitHub Container Registry (GHCR).

---

### Possible Future Improvements

- Support environment variables tobe able to set parameters in the Docker image.
- Improve the normalization algorithm to properly handle single-quoted words and support additional languages.
- Extend support beyond Engadget essays by implementing the `Task` interface for different data sources and resource types.
- Expose metrics (including progress tracking) for monitoring and observability.
- Add the ability to pause and resume work, allowing safe interruption and continuation.

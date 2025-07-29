# Description

A Go application that counts top words based on a list of essays.

---

## 🛠 Requirements

- Go 1.24+ installed
- (Optional) Docker, if you want to run it in a container

## Considerations
- additional to the requirements, words are lowercased:
    - `the` and `The` are the same word
- all pend of word punctation character are removed, except `'`:
    - `end` and `end.` are the same word
    - `its` and `it's` are different words
    - `e.nd` stays `e.nd`

---

## 🚀 Running the App

### Run Locally

```bash
go run ./cmd/app
```



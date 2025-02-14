# GitLab Pipeline Simulator 🚀

[![CI Status](https://github.com/sanchezpaco/gitlab-pipeline-simulator/actions/workflows/ci.yml/badge.svg)](https://github.com/sanchezpaco/gitlab-pipeline-simulator/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sanchezpaco/gitlab-pipeline-simulator)](https://goreportcard.com/report/github.com/sanchezpaco/gitlab-pipeline-simulator)
[![Release](https://img.shields.io/github/v/release/sanchezpaco/gitlab-pipeline-simulator)](https://github.com/sanchezpaco/gitlab-pipeline-simulator/releases)

A powerful CLI tool that predicts which GitLab CI jobs will run under different conditions. Perfect for testing pipeline rules without pushing commits!

## Features ✨

- **Rules Simulation** - Test `if`/`when` conditions locally
- **YAML Expansion** - Resolve anchors, aliases, and `!reference`
- **Regex Support** - Validate branch/tag patterns (=~ operator)
- **Script Preview** - See what commands would execute

## Quick Demo 🎬

```bash
# See what would run for a master branch commit
$ ./glps .gitlab-ci.yml CI_COMMIT_BRANCH=master

🚀 Stage: build
   ✅ Job: build_production
      ├─ Condition: $CI_COMMIT_BRANCH == "master"
      └─ When: always

🚀 Stage: test
   ✅ Job: unit_tests
      ├─ Condition: $CI_PIPELINE_SOURCE == "push"
      └─ When: on_success
```

## Installation 📦

### Binary Downloads
Get pre-built binaries from  [Releases page](https://github.com/sanchezpaco/gitlab-pipeline-simulator/releases)

### From Source

## Releases

```bash
git clone https://github.com/sanchezpaco/gitlab-pipeline-simulator
cd gitlab-pipeline-simulator
make build 
``` 

## Usage Examples 💡

1. Simulate Feature Branch Pipeline

```bash
./glps ci.yml CI_COMMIT_BRANCH=feat/new-auth CI_DEBUG_TRACE=true

🚀 Stage: quality
   ✅ Job: eslint
      ├─ Condition: $CI_COMMIT_BRANCH =~ /^feat\/.*/
      └─ When: manual
   ✅ Job: unittest
      ├─ Condition: $CI_DEBUG_TRACE == "true"
      └─ When: always
```

2. Validate Experimental Tag

```bash
./glps ci.yml CI_COMMIT_TAG=v1.5.0-experimental

🚀 Stage: test
   ✅ Job: chaos_test
      ├─ Condition: $CI_COMMIT_TAG =~ /-experimental$/
      └─ When: manual
```


## CLI Options ⚙️

| Flag | Description |
|-----------|-----------------------------------|
| `-show-scripts`    | Display scripts for matching jobs               |
| `-expand-only`    | Output expanded YAML without simulating a pipeline |
| `-h -help`    | Prints CLI options and examples |


## Current Limitations ⚠️
Not yet implemented (PRs welcome!):
* `include` statements
* `only/except` rules
* Child/parent pipelines

## Development 👩💻

```bash
make test       # Run all tests
make build      # Generates the binary 
```

## Contributing 🤝
Contributions are welcomed! Please check: 

[Contribution guide](CONTRIBUTING.md)

Star ⭐ the repo if you find this useful!
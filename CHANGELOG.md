# [4.0.0](https://github.com/jsburckhardt/gic/compare/v3.0.0...v4.0.0) (2024-11-20)


### Features

* **cmd, config, llm, makefile:** update package documentation and improve formatting and linting ([d5e3753](https://github.com/jsburckhardt/gic/commit/d5e3753ef6d118de28cdd1a7b0b5b54d08e1241d))
* **config:** add new configuration for gic and update instructions in README ([1018b85](https://github.com/jsburckhardt/gic/commit/1018b858848bcdddc839129d51aeac4dc6ce52de))
* **config:** refactor configuration management and environment loading ([1d95ff1](https://github.com/jsburckhardt/gic/commit/1d95ff10afe1af46702a1af40e78bc96f2bdbd77))
* **config:** update README and sample config for environment variables setup ([ed031f6](https://github.com/jsburckhardt/gic/commit/ed031f64f806cb171db9224c1825a3921aa77b9f))
* **logger:** improve logging for error and warning messages ([c8e950a](https://github.com/jsburckhardt/gic/commit/c8e950acd2312cf8f1dce0dc38b106b2e70d9722))


### BREAKING CHANGES

* **config:** This refactor changes the structure of the configuration management and requires updating the environment variable setup in the deployment.

# [3.0.0](https://github.com/jsburckhardt/gic/compare/v2.4.0...v3.0.0) (2024-10-24)


### Features

* **cmd, git:** add support for generating commit messages based on the main branch for pull requests ([3c1e8bb](https://github.com/jsburckhardt/gic/commit/3c1e8bbc087cb31cb3466da5e742eb5fcec8cd19))
* **git:** enhance diff handling based on configuration ([836c7dd](https://github.com/jsburckhardt/gic/commit/836c7dd85880ce95d049d3f24fe3f0955b1a0d6a))


### BREAKING CHANGES

* **git:** Changes the `Commit` function by removing the `pr` parameter. Now, the function relies solely on the `ShouldCommit` and `PR` fields in the configuration.

# [2.4.0](https://github.com/jsburckhardt/gic/compare/v2.3.0...v2.4.0) (2024-09-13)


### Features

* **docs, config, devcontainer:** enhance README and configuration for gic usage ([7bec5ab](https://github.com/jsburckhardt/gic/commit/7bec5ab4d38b94337dbbca79299de051603bbb83)), closes [#12345](https://github.com/jsburckhardt/gic/issues/12345)

# [2.3.0](https://github.com/jsburckhardt/gic/compare/v2.2.0...v2.3.0) (2024-09-13)


### Features

* **config:** add create-sample-config flag and sample configuration generation ([2d34d83](https://github.com/jsburckhardt/gic/commit/2d34d83cc4757336d0131a5d1cae35d23dff017c))

# [2.2.0](https://github.com/jsburckhardt/gic/compare/v2.1.0...v2.2.0) (2024-09-12)


### Bug Fixes

* **cmd:** correct string concatenation for commit message logging ([2ef6146](https://github.com/jsburckhardt/gic/commit/2ef61469e14bc4ffdf8a286dc0c1f1234f08ec82))
* **cmd:** standardize log output for commit message formatting ([067b904](https://github.com/jsburckhardt/gic/commit/067b904c449b6843f5e953fe0dc83d747f91074f))


### Features

* **devcontainer, ci:** update devcontainer configuration and CI script ([5222f60](https://github.com/jsburckhardt/gic/commit/5222f60520140875506d8663551c72f02afe4068))
* **logging:** add info log for commit messages before execution ([5ce4d61](https://github.com/jsburckhardt/gic/commit/5ce4d61b11eaf3c3154edd8f1180e557a6425531))

# [2.1.0](https://github.com/jsburckhardt/gic/compare/v2.0.0...v2.1.0) (2024-09-11)


### Bug Fixes

* **git:** log commit message before committing changes ([28c0e56](https://github.com/jsburckhardt/gic/commit/28c0e563905f4b6bd41bc1fec30033b29b4d0736))


### Features

* **logger:** enhance logging mechanism with caller info and message formatting ([ca6d629](https://github.com/jsburckhardt/gic/commit/ca6d629a0529cde6f76c30ab16d5a204513e5f2e))
* **logger:** enhance logging with caller information and control over source display ([cf84973](https://github.com/jsburckhardt/gic/commit/cf8497337a029783158cd20f7a8f47f797c8e747))
* **logger:** implement structured logging with adjustable log levels ([7ea38b1](https://github.com/jsburckhardt/gic/commit/7ea38b18156c49a314328b84e99e9a23f6d02e2b))

# [2.0.0](https://github.com/jsburckhardt/gic/compare/v1.1.0...v2.0.0) (2024-09-02)


### Features

* **cmd:** refactor command execution and configuration handling ([272ffde](https://github.com/jsburckhardt/gic/commit/272ffde4bbe2f087a88c567f139446f51fa52469))
* **git:** implement commit functionality with message suggestion ([c91d03c](https://github.com/jsburckhardt/gic/commit/c91d03cd94e8e0ef89a8288a453b84a14f052e6b))
* **llm, api:** integrate Ollama API for generating commit messages ([9f228cb](https://github.com/jsburckhardt/gic/commit/9f228cb05c187702b7c9d0ced88d01acc4a4e877))


### BREAKING CHANGES

* **cmd:** The `Commit` parameter in the `git.Commit` function has been changed to accept the entire `config.Config` struct instead.
* **llm, api:** The connection type must now include support for "ollama".

# [1.1.0](https://github.com/jsburckhardt/gic/compare/v1.0.0...v1.1.0) (2024-08-29)


### Features

* **init:** add main package and command execution ([220b9d4](https://github.com/jsburckhardt/gic/commit/220b9d4f60af4de86cd09f82a348079a0437d8c0))

# 1.0.0 (2024-08-23)


### Features

* **ci:** add GoReleaser configuration and integrate with semantic-release ([7ad2eca](https://github.com/jsburckhardt/gic/commit/7ad2eca17699c3f892c586e8a1e8fcb5978eecb3))

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

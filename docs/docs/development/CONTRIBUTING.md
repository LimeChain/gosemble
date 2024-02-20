# Contributing to Gosemble

Thank you for considering contributing to Gosemble! We welcome contributions from everyone. To maintain a high standard of quality, please follow these guidelines when contributing to the project.

## Signing Your Commits

To contribute to Gosemble, you must have a verified GPG key added to your GitHub account. This ensures the integrity and authenticity of your contributions. For more information on adding a GPG key to your GitHub account, refer to [GitHub's official documentation](https://docs.github.com/en/github/authenticating-to-github/managing-commit-signature-verification/adding-a-new-gpg-key-to-your-github-account).

## Commit Message Guidelines

We follow the [Conventional Commits specification](https://www.conventionalcommits.org/) for our commit messages. Each commit message should have a clear and structured format that conveys the type, scope, and description of the changes made. This helps in understanding the purpose of each commit and automating tasks like versioning and changelog generation.

## Merge Process

When submitting a new pull request (PR), please ensure the following:

- **Target the `develop` branch**: All new PRs must target the `develop` branch.
- **Use our PR template**: Utilize our [PR template](../../../.github/PULL_REQUEST_TEMPLATE.md) to ensure all necessary information is included.
- **Include updated runtime.wasm**: Our test workflow relies on having an actual version of the [runtime.wasm](../../../build/runtime.wasm) file. you should [build](./build.md) the runtime and commit the changes.
- **Pass the tests**: Ensure that all tests are passing. [Read more](./test.md).
- **Pass the coverage workflow**: Ensure that your PR is passing the [test coverage workflow](../../../.github/workflows/coverage.yaml).
- **Update the docs**: Ensure that all relevant changes are reflected in the docs.
- **Follow our style guide**: Check our [style guide](#style-guide) and maintain consistent coding style.
- **Have 2 approvals from the core dev team**: PRs should have approvals from at least two members of the core development team before merging.

Additional recommendations:
- Write integration tests for logic related to pallets. [See example test](../../../runtime/balances_set_balance_test.go).
- Write benchmark tests for new extrinsic calls. [Read more](./benchmarking.md).
- Check if you can successfully start a substrate node with the runtime. [See guide](../tutorials/start-a-network.md).

## Release Process

We use the `develop` branch for ongoing development. Once the development is complete and ready for release, changes are merged into the `master` branch.

When changes are merged into the `master` branch, a release process is triggered automatically through our [CI workflow](../../../.github/workflows/release.yaml). We use a tool called [release-please](https://github.com/marketplace/actions/release-please-action) that generates a pull request with the changelog and creates a release on GitHub. This helps in automating the release process and maintaining an organized changelog for each release.

Currently, we aim to have a release after each phase. For more information about milestones and phases, check our [project page](https://github.com/orgs/LimeChain/projects/5/views/8).

## Style Guide

We try to follow common practices described in [Effective Go](https://go.dev/doc/effective_go) and [Google Style Guide for Go](https://google.github.io/styleguide/go/) to maintain consistency across the codebase. Following a consistent coding style makes the code more readable and maintainable for everyone.

- **Custom errors**: New polkadot-related custom errors should implement the error interface. See [#271](https://github.com/LimeChain/gosemble/issues/271) and linked PRs.
- **Error handling**: We use critical logging for resolving errors(panicing). We add logger to modules through dependency injection and we only resolve errors(critical log) in the api modules. See [#315](https://github.com/LimeChain/gosemble/pull/315).
- **SCALE types**: Certain [SCALE types](../overview/runtime-architecture.md#scale-codec) like the Result types are common in Rust which the original Substrate implementation is based on, but add unneeded complexity in Go. When we have these types as return types it's preferred to encode the data into the required type as late as possible - just when you need to return it in the api module, instead of building your logic around them. See [#322](https://github.com/LimeChain/gosemble/pull/322).
---
layout: default
title: Contributing guidelines
permalink: /contributing/contributing
---

# Contributing

Thank you for considering contributing to Gosemble! We welcome contributions from everyone. To maintain a high standard of quality, please follow these guidelines when contributing to the project.

## Signing Your Commits

To contribute to Gosemble, you must have a verified GPG key added to your GitHub account. This ensures the integrity and authenticity of your contributions. For more information on adding a GPG key to your GitHub account, refer to [GitHub's official documentation](https://docs.github.com/en/github/authenticating-to-github/managing-commit-signature-verification/adding-a-new-gpg-key-to-your-github-account).

## Commit Message Guidelines

We follow the [Conventional Commits specification](https://www.conventionalcommits.org/) for our commit messages. Each commit message should have a clear and structured format that conveys the type, scope, and description of the changes made. This helps in understanding the purpose of each commit and automating tasks like versioning and changelog generation.

## Merge Process

When submitting a new pull request (PR), please ensure the following:

- **Target the `develop` branch**: All new PRs must target the `develop` branch.
- **Use our PR template**: Utilize our [PR template](https://github.com/limechain/gosemble/tree/develop/.github/PULL_REQUEST_TEMPLATE.md) to ensure all necessary information is included.
- **Include updated runtime.wasm**: Our test workflow relies on having an actual version of the [runtime.wasm](https://github.com/limechain/gosemble/tree/develop/build/runtime.wasm) file. you should [build](../development/build.md) the runtime and commit the changes.
- **Pass the tests**: Ensure that all tests are passing. [Read more](../development/test.md).
- **Pass the coverage workflow**: Ensure that your PR is passing the [test coverage workflow](https://github.com/limechain/gosemble/tree/develop/.github/workflows/coverage.yaml).
- **Update the docs**: Ensure that all relevant changes are reflected in the docs.
- **Follow our style guide**: Check our [style guide](./style-guide.md) and maintain consistent coding style.
- **Have 2 approvals from the core dev team**: PRs should have approvals from at least two members of the core development team before merging.

Additional recommendations:
- Write integration tests for logic related to pallets. [See example test](https://github.com/limechain/gosemble/tree/develop/runtime/balances_set_balance_test.go).
- Write benchmark tests for new extrinsic calls. [Read more](../development/benchmarking.md).
- Check if you can successfully start a substrate node with the runtime. [See guide](../tutorials/start-a-network.md).

## Release Process

We use the `develop` branch for ongoing development. Once the development is complete and ready for release, changes are merged into the `master` branch.

When changes are merged into the `master` branch, a release process is triggered automatically through our [CI workflow](https://github.com/limechain/gosemble/tree/develop/.github/workflows/release.yaml). We use a tool called [release-please](https://github.com/marketplace/actions/release-please-action) that generates a pull request with the changelog and creates a release on GitHub. This helps in automating the release process and maintaining an organized changelog for each release.

Currently, we aim to have a release after each phase. For more information about milestones and phases, check our [project page](https://github.com/orgs/LimeChain/projects/5/views/8).

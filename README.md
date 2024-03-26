# Semver Version Bump

This is a GitHub action that finds any labels on a PR associated with a merged commit that indicate whether the PR necessitates a version bump. If it does, it also determines if the version bump should be a patch, minor, or major bump. If no labels are found, it defaults to the bump configured with the `DEFAULT_BUMP` input.

## Inputs

These should be set as environment variables in the workflow file.

| Field | Required | Default | Description |
| --- | --- | --- | --- |
| GH_TOKEN | true | | The GitHub token used to authenticate with the GitHub API. |
| PATCH_LABELS | false | patch | A list of labels that indicate a patch bump, comma-separated. |
| MINOR_LABELS | false | minor | A list of labels that indicate a minor bump, comma-separated. |
| MAJOR_LABELS | false | major | A list of labels that indicate a major bump, comma-separated. |
| NO_BUMP_LABELS | false | no bump | The labels that indicate no bump. |
| DEFAULT_BUMP | false | none | The default bump if no labels are found. |

> [!NOTE]
> If labels corresponding to multiple bump types are found, the bump type with the highest precedence will be used. The precedence is as follows:
> 1. `no bump`
> 1. `major`
> 1. `minor`
> 1. `patch`

## Outputs

| Field | Description |
| --- | --- |
| bump | The bump type. Will be one of: `none`, `major`, `minor`, or `patch`. |

## Sample Usage

```yaml
name: Version Bump
on:
  push:
    branches:
      - main

jobs:
  version-bump:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: read
    steps:
      - name: Checkout
        uses: actions/checkout@v2

      - name: Version Bump
        id: version-bump
        uses: heru-inc/semver-version-bump@latest
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          DEFAULT_BUMP: patch
          # ... #

      - name: Print Bump
          run: echo "Bump type: ${{ steps.version-bump.outputs.bump }}"
```

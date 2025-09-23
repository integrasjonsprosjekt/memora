# Contributing

Before contributing, please familiarise yourself with the technologies and frameworks listed below.
This helps us maintain a consistent structure, and lowers the overhead of context switching.

## Development

We use Nix to manage development dependencies/environments.
As a companion to Nix, we suggest using direnv locally for better integration with editors, and automatic environment activation.

In addition to Nix, we also have Docker containers configured for both development and production use.
When using the development shell exposed by the Nix flake, the Docker Compose profile is already set as `dev`.
If you're not using Nix, set the profile manually with `export COMPOSE_PROFILES=dev`.

### Linting and formatting

The Nix flake makes all required tools available from the development shell.

All code has to pass linting and formatting checks to be merged into the default branch.

### Git hooks

We use [husky](https://github.com/typicode/husky) to manage our Git hooks. This ensures the proper linting, formatting, and convention checks are run before pushing code to a remote.

If you're using the Nix flake, husky is automatically initialized for you.\
Otherwise, you will have to run `npx husky` once to properly set up husky locally.

## Naming conventions

### Commits

Commits should follow the _conventional commits_ specification.
Please read the following cheat sheet for a simplified specification: https://gist.github.com/qoomon/5dfcdf8eec66a051ecd85625518cfd13

### Branches

Branches should be named using a very simplified version of _gitflow_.
The branch types/prefixes mostly follow the same convention as the commit types, with the exception of `feat` and `fix`, which are called `feature` and `bugfix`/`hotfix` respectively.

> [!NOTE]
> - **Do** prefix the branch name with one of the types `feature`, `bugfix`, `hotfix`, `chore`, `docs`, `ci`, `build`, `test`, `refactor`, or `style`:\
>   <code><b>feature/</b>description</code>, <code><b>hotfix/</b>description</code>, etc., **not** <code><b>feat</b>/description</code> or <code><b>fix</b>/description</code>
> - **Do not** use scopes:\
>   <code>type/description</code>, **not** <code>type<b>(scope)</b>/description</code>
> - **Do** separate the type and description with a forward slash:\
>   <code>type<b>/</b>description</code>, **not** <code>type<b>:</b> description</code> or <code>type<b>!</b>/description</code>
> - **Do not** capitalise any part of the name (including the first character):\
>   <code>feature/add-markdown-support</code>, **not** <code>feature/<b>A</b>dd-markdown-support</code>
> - **Do** use the imperative, present tense — "add", **not** "added" or "adds":\
>   <code>feature/<b>add</b>-markdown-support</code>, **not** <code>feature/<b>added</b>-markdown-support</code>
> - **Do** use hyphens to separate words:\
>   <code>feature/add<b>-</b>markdown<b>-</b>support</code>, **not** <code>feature/add markdown support</code>
> - **Do** try to keep the description short and descriptive.

### Pull Requests

Pull requests should be named the same as the branch.
Since pull requests in GitHub behave similarly to issues, this makes it easier to know what you're looking at.

### Issues

Issue names should be prefixed with the type of issue, followed by a brief description of the contents, e.g. `[Feature] Markdown support`.


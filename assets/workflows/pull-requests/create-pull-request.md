---
uri: workflows://create-pull-request
name: Create Pull Request Workflow
description: "Use when a committed branch is ready to open a pull request. Covers local checks, opening a ready pull request, watching CI, and handling comments. Does not merge."
languages:
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - uri: standards://git/pull-requests
      when: "writing the title or description"
    - uri: standards://git/commit-messages
      when: "a later fix needs a commit message"
    - uri: standards://git/commit-staging
      when: "a later fix must be split into commits"
    - uri: standards://git/operations
      when: "pushing or rewriting history after the pull request exists"
    - uri: standards://git/branches
      when: "the branch must be created or renamed before opening"
    - uri: workflows://investigate-ci-failures
      when: "a pull request check failed, was cancelled, or timed out"
---

# Create Pull Request Workflow
Use when a committed branch is ready to open a pull request. Covers local checks, opening a ready pull request, watching CI, and handling comments. Does not merge.

Read this file. Read a related resource only when its `when` condition matches this step.

Use this when the branch already contains the commits that should be reviewed. Follow `standards://git/pull-requests` when writing the title and description. Read `standards://git/operations`, `standards://git/commit-messages`, or `standards://git/commit-staging` only when a later fix commit or push matches that resource's `when`.

A pull request here includes the host's equivalent, such as a GitLab merge request. The title, description, ready state, CI, and comment rules are the same. The commands are not.

- ALWAYS assume the work is already committed. Do not commit, rename the branch, or clean the tree in order to start validation.
- ALWAYS open a ready-for-review pull request once local checks pass. Do not ask for a second confirmation. Do not open a draft.
- ALWAYS repeat CI investigation and comment handling, with no cycle cap, until the completion rule in Step 6 is met.
- NEVER merge, enable auto-merge, or approve the pull request.
- Follow `standards://git/operations` for every push and history change. Do not restate those rules here.

## Host tooling

Identify the host from the remote URL before choosing a tool.

```bash
git remote get-url origin
```

Use the branch upstream remote when it is not `origin`.

- ALWAYS use that host's CLI or API to create, edit, watch checks, list comments, reply, and resolve.
- NEVER point `gh` or the GitHub API at a remote that is not GitHub or GitHub Enterprise.
- GitHub.com and GitHub Enterprise: `gh`, or the GitHub REST or GraphQL API. For Enterprise, target that hostname. Do not assume `github.com`.
- GitLab.com and self-managed GitLab: `glab`, or the GitLab API.
- Bitbucket Cloud: the Bitbucket REST API. Do not assume a GitHub-style CLI. Bitbucket Data Center uses that server's API.
- Gitea and Forgejo: `tea`, or that server's API.
- Gogs: that server's API, or the CI system attached to the repository. Gogs has no assumed CI CLI.
- These hosts are examples, not a complete list. Another host uses the CLI or API it documents.
- IF the host cannot be identified, or its tool is missing, unauthenticated, or unknown, STOP and ask. Do not invent commands.
- Command blocks below are GitHub examples of one option. On another host, perform the same action with that host's CLI or API.

## Step 1: Run local validation

Discover the checks this repository already uses. Read its agent rules, `Makefile`, package scripts, and CI configuration. Do not assume `make`.

- ALWAYS run the repository's tests, lint, format check, typecheck, vet, and security or vulnerability checks when the repository defines them.
- ALWAYS treat a non-zero exit as a failure.
- IF a failure already exists on the base branch, and this branch did not add it, record that and continue.
- IF you cannot prove a failure is pre-existing, treat it as introduced by this branch and STOP. Do not open the pull request.
- NEVER check out the base branch in this working tree to compare. Use the base branch's CI, or a separate worktree, when a comparison is required.
- Do not commit in this step. A dirty tree is allowed until the push in Step 4.

## Step 2: Reproduce CI locally

Before creating the pull request, read the CI definitions that will run when it opens. Discover them from the repository. Do not assume one host's layout.

- Example layouts, not a required path: `.github/workflows/` for GitHub and GitHub Enterprise, `.gitlab-ci.yml` for GitLab, `bitbucket-pipelines.yml` for Bitbucket, `.gitea/workflows/` for Gitea, and a `Jenkinsfile` or other project-specific pipeline file. Gitea may also read `.github/workflows/`. Use the file the repository actually runs.
- Run CI commands that Step 1 did not already run.
- PREFER the repository's local entrypoint when a workflow calls it, rather than retyping every flag.
- IF a step needs secrets, a hosted runner, or an external environment, list it and leave it for CI. Do not record it as passed.
- IF a locally reproduced step fails, fix nothing in this step unless the user has already asked for that fix. STOP and report the failure. Do not open the pull request.

## Step 3: Generate the title and description

Follow `standards://git/pull-requests` for the title and the description. Do not restate that format here.

- ALWAYS read every commit subject on the branch, not only the latest, and build the title from those subjects:

```bash
git log --format='%s' <base>..HEAD
```

- Do not let the host CLI fill the title or body from commit text, and do not open an editor or browser to write them. Example, GitHub: do not pass `--fill`, `--fill-first`, `--fill-verbose`, `--editor`, or `--web`.
- IF a commit subject does not follow `standards://git/commit-messages`, STOP and ask. Follow `standards://git/operations`. Do not rewrite that commit.
- The base branch is the existing pull request base, otherwise the repository default branch. IF that cannot be determined, STOP and ask.

## Step 4: Push and open the pull request

- IF `git status --porcelain` is not empty, STOP. Do not commit, stash, or discard changes in order to open the pull request.
- ALWAYS push with git before creating the pull request. Push does not depend on the host CLI.
- IF the branch has an upstream, `git push` to that upstream.
- IF it has no upstream, `git push -u origin HEAD`.
- IF the push is rejected, STOP and ask. Follow `standards://git/operations`.
- IF a pull request for this branch already exists, do not create another. Push, then update its title and body with the host CLI. IF it is a draft, mark it ready.
- IF no pull request exists, create a ready one with that host's CLI or API. Example, GitHub:

```bash
gh pr create --base <base> --title "<title>" --body-file <body-file>
```

- Example, GitHub: do not pass `--draft`.
- Do not add reviewers, labels, or assignees unless the user asked.

## Step 5: Watch CI until every check finishes

- ALWAYS watch until every check has finished. Do not stop at the first failure. Example, GitHub: `gh pr checks --watch` without `--fail-fast`.

- A failed, cancelled, or timed-out check is a failure. A skipped check is not.
- IF any check fails, follow `workflows://investigate-ci-failures`. That workflow tries solutions in a loop and does not push. IF the user interrupts or aborts, stop. Do not call it again.
- After it returns, commit the fix if the working tree is dirty. Follow `standards://git/commit-messages`, `standards://git/commit-staging`, and `standards://git/operations`.
- IF the fix is already committed, do not create another commit.
- ALWAYS re-run Step 1 before pushing a fix. The tree must be clean at that push.
- ALWAYS push, then update the title and description if the commits changed them. Follow `standards://git/pull-requests`.
- ALWAYS return to this step after that push. Repeat until every check has succeeded, or the user interrupts or aborts. A red check after that watch is the same failure loop, not a new task.

## Step 6: Review every comment

WHEN every CI check has succeeded, review every comment on the pull request. Include humans and bots. Examples of bots: Bugbot, Wiz, CodeRabbit. The list is not exhaustive.

A command that lists only timeline comments is not enough. On GitHub, `gh pr view --comments` is that command.

- ALWAYS also list inline review comments and whether each thread is resolved.
- Example, GitHub, list inline comments with:

```bash
gh api repos/{owner}/{repo}/pulls/<number>/comments
```

- Example, GitHub: read review threads, including `isResolved`, with `gh api graphql`. A thread id is required later to resolve a thread that a code fix addressed. Another host uses its own thread id.
- IF there are no comments, or every thread is already answered and the resolution rule below is already met, the workflow is complete.
- IF any thread is unanswered, or a real issue is not yet fixed, handle each thread:

1. Read the finding, the code around it, the plan for that behavior, and the conversation.
2. Decide whether it is a question or an issue.
3. IF it is a question, reply in the thread and leave it unresolved. Cite the code, documents, rules, skills, or MCP resources that support the answer.
4. IF it is an issue, reproduce it where possible, or review the code it describes. Do not call it invalid in order to avoid the fix.
5. IF the issue is real, implement the fix. Commit and push as in Step 5. Resolve only that thread, and only after the fixing commit is on the remote. Then return to Step 5.
6. IF the issue is not real, reply in the thread with a detailed claim and the same kind of evidence. Leave the thread unresolved.

Reply inside the review thread. On GitHub, `gh pr comment` creates an issue comment and does not answer an inline thread.

- On GitHub, the reply target must be the top-level review comment id. Replies to replies are not supported.

```bash
gh api --method POST repos/{owner}/{repo}/pulls/<number>/comments/<comment_id>/replies -f body='...'
```

- On GitHub, resolve a thread only with the GraphQL `resolveReviewThread` mutation, and only when this workflow's code fix addressed that thread.
- On another host, reply and resolve with that host's API. IF the host cannot resolve a thread, reply and leave it unresolved. Do not send a GitHub mutation.
- NEVER resolve a question.
- NEVER resolve a thread whose finding was rejected.
- AFTER new comments appear, including comments from the next CI run, handle them the same way. Repeat Step 5 and this step until the completion rule is met.

The workflow is complete when all of the following are true:

- Every CI check has succeeded.
- Every comment thread has a reply.
- Threads fixed in code are resolved, or the host has no resolve action and the reply says so.
- Questions and rejected findings are still unresolved.
- The pull request is not merged.

## Checklist

- [ ] Host was identified from the remote, and that host's tool was used
- [ ] Local validation passed, or every remaining failure was proven to exist on the base branch
- [ ] CI commands that can run locally were re-run before the pull request was opened
- [ ] Title and description follow `standards://git/pull-requests`
- [ ] Push and history follow `standards://git/operations`
- [ ] Pull request is ready for review, not a draft
- [ ] Every CI check succeeded
- [ ] Every thread has a reply
- [ ] Only threads fixed in code are resolved
- [ ] Pull request was not merged, auto-merged, or self-approved

---
uri: workflows://investigate-ci-failures
name: Investigate CI Failures Workflow
description: "Use when a pull request check failed, was cancelled, or timed out. Repeat until that check is green or the user aborts. Does not push or merge."
languages:
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - uri: workflows://create-pull-request
      when: "host tooling is needed, or the fix is ready to hand back to the open pull request"
    - uri: standards://git/operations
      when: "the fix requires a push or history rewrite"
    - uri: standards://git/commit-messages
      when: "the fix needs a commit message"
    - uri: standards://git/commit-staging
      when: "the fix must be split into commits"
    - uri: standards://agentic-programming/architecture
      when: "the repository's docs/ or .agents/ layout is needed to find the checks"
---

# Investigate CI Failures Workflow
Use when a pull request check failed, was cancelled, or timed out. Repeat until that check is green or the user aborts. Does not push or merge.

Read this file. Read a related resource only when its `when` condition matches this step.

Use this when a pull request check has failed, been cancelled, or timed out. `workflows://create-pull-request` Step 5 calls this workflow for that case, then watches CI and calls it again if the check is still red. The same steps apply when this workflow is started on its own.

Follow the host-tooling rules in `workflows://create-pull-request`. The calls in Examples are options, not the procedure. GitHub is one option.

- ALWAYS repeat the solution loop until the recorded CI check is green, or the user interrupts or aborts.
- IF the user interrupts or aborts, stop immediately. Do not start another attempt.
- ALWAYS investigate every finished failure from this run before returning. Do not stop at the first red check.
- ALWAYS keep the original error text. Do not replace it with a paraphrase.
- NEVER push, merge, or approve from this workflow. History changes follow `standards://git/operations`.
- NEVER re-run the remote workflow as the fix. A retry is not a solution.
- NEVER skip, disable, or weaken a check to make it pass.
- IF the host tool is missing, unauthenticated, or unknown, STOP and ask. Do not invent commands.

## Examples

Use the host's CLI or its API. Either can satisfy Steps 1 and 2. Do not run a GitHub command against another host.

GitHub and GitHub Enterprise, CLI:

```bash
gh pr checks
gh run list --branch "$(git branch --show-current)" --commit "$(git rev-parse HEAD)" --status failure
gh run view <run-id> --verbose
gh run view <run-id> --log-failed
```

GitHub and GitHub Enterprise, API: list check runs for the head SHA, then download the failed job log. For Enterprise, target that hostname. Do not assume `github.com`. When checks are still running, `gh pr checks --watch` without `--fail-fast` is one way to wait.

GitLab, CLI:

```bash
glab ci status
glab ci trace <job-id>
```

GitLab, API: list the pipeline's jobs, then fetch that job's trace. `glab ci status --wait` is one way to wait until the pipeline finishes. This applies to GitLab.com and a self-managed GitLab.

Bitbucket Cloud, API: list pipelines for the repository, then fetch the step log. There is no assumed CLI.

```text
GET /2.0/repositories/{workspace}/{repo_slug}/pipelines
GET /2.0/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/steps/{step_uuid}/log
```

The UUID is the pipeline or step UUID, not the display name. Bitbucket Data Center uses that server's API, not the Cloud paths above.

Gitea, CLI. Forgejo can use the same CLI when `tea` is logged in to that server:

```bash
tea actions runs list --status failure --branch "$(git branch --show-current)"
tea actions runs view <run-id>
tea actions runs logs <run-id> --job <job-id>
```

Gitea, API: list Actions runs for the repository, then fetch the run or job log from that server.

Gogs has no assumed CI CLI and no GitHub-style check log. Use the CI system attached to that repository and its API. IF that system is unknown, STOP and ask.

These hosts are examples, not a complete list. IF the remote is a host not shown here, use the CLI or API that host documents. IF neither is known, STOP and ask.

## Step 1: Record the failing checks and steps

Read the pull request's CI statuses and messages with the host's CLI or API. Record every check that failed, was cancelled, or timed out. A skipped check is not a failure.

For each failure, record the check name, conclusion, workflow name, job name, step name, run URL, and head SHA. The step name is the step inside the job, not only the job name.

- IF checks are still pending, wait until they finish. Do not stop at the first failure.
- IF a run cannot be matched to a step, record that gap. Do not guess the step name.
- Group failures that share one command and one error. Keep unrelated failures separate.

## Step 2: Isolate the error from the logs

Read the logs of each failed step with the host's CLI or API. Locate the exact error message, stack trace, and error code, and the command that produced them.

- ALWAYS quote the failing lines, including the error code and the first frame of the stack that points at this repository.
- ALWAYS record the command the step ran, and the exit status when the log shows one.
- IF the log is missing, truncated, or has no error line, record that. Do not invent an error message.
- A cancelled or timed-out run may have no assertion failure. Record the conclusion and the last step that started.
- IF the failure is a missing secret, a permission, or a hosted service this environment cannot access, STOP and ask. Do not work around it in product code.

## Step 3: Read what the failing step is supposed to do

Open the CI definition for the failing step. Example layouts are in `workflows://create-pull-request` Step 2. Use the file the repository actually runs. Do not assume one of those paths exists.

- ALWAYS read the job and step that failed, including the command, its inputs, and the condition that decides it passed.
- ALWAYS note whether that command can be run locally. A step that needs secrets or a hosted runner cannot be fully reproduced here.
- Distinguish a product failure, a test failure, a workflow-wiring failure, and an environment failure. The later steps must use that distinction.
- IF the workflow file and the log name different commands, trust the log for what ran and record the mismatch.

## Step 4: Compare the failure with the intended result

Search the conversation and the repository for plans, designs, and decisions that say what this change was supposed to do. Use documents the branch or conversation already names. Read that repository's agent index, such as `AGENTS.md`, for where it keeps this evidence.

- ALWAYS search the repository. Do not stop at one familiar path.
- IF the repository follows `standards://agentic-programming/architecture`, use the plan and decision files that standard names. Those paths are not a fact about every repository.
- Other example locations, not a required layout: `plans/`, `plan.md`, `docs/design/`, `docs/adr/`, and `adr/`. A repository may use none of these.
- ALWAYS treat the documents you find as evidence of intent, not as a reason to ignore the log.
- IF a plan and the failing check disagree, record both. Do not silently pick one.
- IF no plan exists, continue from the log, the workflow step, and the code. Do not invent a plan.
- The intended result is the behavior the check was written to protect, unless a plan explicitly changes that behavior and the check has not been updated yet.

## Step 5: Propose solutions and test them against the evidence

Propose more than one way to fix the recorded failure. Do not edit code in this step.

- ALWAYS propose at least two distinct solutions when more than one cause fits the log.
- IF the log, the workflow step, and the plan pin a single cause, name that cause and at least one rejected alternative. Say which recorded line rules the alternative out.
- ALWAYS test each proposal against the quoted error, the stack, the command the step ran, and the intended result from Step 4. A proposal is effective only when it would change the condition that produced that error.
- NEVER accept a proposal that skips the check, deletes the assertion, or retries the job without changing the cause.
- NEVER accept a proposal that contradicts a plan unless the log shows the plan itself is what failed. IF that conflict changes product behavior, STOP and ask.
- Leave the rejected proposals recorded, with the evidence that rejected them.
- Keep every surviving proposal. The solution loop tries them one at a time. Do not discard an untried proposal because another one was implemented.

## Step 6: Implement the most effective fix

Choose the proposal that the evidence supports. Implement that fix only.

- IF two surviving proposals are equally supported, choose the smaller change that fixes the recorded cause.
- IF no proposal survives Step 5, STOP and ask. Do not guess.
- ALWAYS implement every unrelated failure group recorded in Step 1 before returning, unless a group was stopped in Step 2 or Step 5.
- Keep the edits separable under `standards://git/commit-staging`.
- IF the failing command can run locally, re-run that command. A local pass does not prove a hosted-only step. Record that limit.
- Do not treat this attempt as the end of the workflow. Continue at the solution loop.
- Do not push or watch the remote run here. Commit only in the hand-back rules below.

## Solution loop

After each attempt, decide from the new result. Do not stop because one proposal was tried.

- IF the user interrupts or aborts, stop immediately.
- IF the recorded check is green, the loop is finished. Use the hand-back rules.
- IF the attempt still fails and the new output changes the cause, quote that output and return to Step 2. Propose a new set in Step 5 from that evidence, then try the most effective untried proposal.
- IF the attempt still fails and the cause is unchanged, do not propose the same set again. Try the next most effective untried proposal from the current set.
- IF a failed attempt is still uncommitted, change or remove those edits before the next attempt when they would interfere. Do not commit a failed attempt just to try the next one.
- IF an earlier attempt is already committed, leave that commit. The next fix is a new commit. Follow `standards://git/operations`.
- IF no untried proposal remains, and the latest evidence does not support a new set, STOP and ask. Do not invent a proposal that fails the Step 5 evidence test.
- A local pass is not remote CI green. Hand back so the check can be watched. IF that watch is still red, resume this workflow from Step 1 with the new logs. That is the same loop, not a new task.

## Hand back

This workflow is complete only when the recorded CI check is green, or the user has interrupted or aborted. A hand-back is not completion while the check is still red.

- IF `workflows://create-pull-request` called this workflow, return to its Step 5 after the current attempt is ready to commit. IF the check is still red, that step calls this workflow again.
- IF this workflow was started on its own and a pull request is already open, commit the fix under `standards://git/commit-messages` and `standards://git/commit-staging`, then continue at `workflows://create-pull-request` Step 5. Do not push from here.
- IF no pull request exists, commit the fix under those standards and stop. Open and watch the pull request with `workflows://create-pull-request`.

## Checklist

- [ ] Every failed, cancelled, and timed-out check was recorded with its step
- [ ] The error text, stack, error code, and failing command were quoted, or their absence was recorded
- [ ] The workflow step was read for what should have happened
- [ ] Plans in the conversation and repository were checked, or their absence was recorded
- [ ] More than one solution was compared against that evidence, or a single cause was recorded with a rejected alternative
- [ ] The implemented fix is the one the evidence supported
- [ ] A failed attempt continued the solution loop, or the user aborted
- [ ] The failing local command was re-run when it could be
- [ ] No check was skipped or weakened, and the remote workflow was not re-run as the fix
- [ ] Nothing was pushed or merged, and history follows `standards://git/operations`

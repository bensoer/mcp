---
uri: workflows://skill-authoring
name: Skill Authoring Workflow
description: "Use when creating or revising an Agent Skill from real work. Follow standards://skill-authoring for format and content rules. Do not load this for a single rule check."
related_resources:
    - uri: standards://skill-authoring
      when: "applying a format, description, context, or script rule"
---

# Skill Authoring Workflow
Use when creating or revising an Agent Skill from real work. Follow standards://skill-authoring for format and content rules. Do not load this for a single rule check.

Read this file. Read a related resource only when its `when` condition matches this step.

Follow `standards://skill-authoring` for every format and content rule. Do not restate those rules here.

## Step 1: Extract real expertise

Do not generate a skill from general training knowledge.

1. Prefer a hands-on task. Complete it with an agent, then extract the steps that worked, the corrections you made, the input and output formats, and the project facts the agent did not already know.
2. Otherwise synthesize from project artifacts: runbooks, schemas, review comments, and real failures. Verify the result captures those artifacts, not generic advice.

## Step 2: Draft the skill

1. Create the skill directory and `SKILL.md`.
2. Apply `standards://skill-authoring` for the name, description, body, progressive disclosure, instruction patterns, and scripts.
3. Put only the rules the agent would get wrong. Do not copy this workflow into the skill.

## Step 3: Refine with real execution

1. Run the skill against a real task.
2. Read the execution trace, not only the final output.
3. When you have to correct the agent, add that correction to the gotchas section required by `standards://skill-authoring`.
4. Repeat once. Keep a change only if the next run uses it.

## Step 4: Optimize the description

1. Apply the description rules in `standards://skill-authoring`.
2. Run the optimization loop defined there. Do not invent a second loop here.
3. Revise the general category of a miss. Do not paste keywords from one failed query.

## Step 5: Evaluate output quality

1. Apply the evaluation rules in `standards://skill-authoring`.
2. Run one case with the skill and one case without it, each in a clean context.
3. Stop when a further revision does not change the result.

## Checklist

- [ ] The skill came from a real task or project artifacts
- [ ] Format and content rules were applied from `standards://skill-authoring`, not copied from this file
- [ ] The description says what the skill does and when to use it
- [ ] A with-skill run was compared with a without-skill run

---
uri: workflows://skill-authoring
name: Skill Authoring Workflow
description: Step-by-step workflow for creating an effective Agent Skill — from extracting real expertise to drafting, refining, optimizing, evaluating, and iterating. Follows the standards defined in standards://skill-authoring.
related_resources:
    - standards://skill-authoring
---

# Skill Authoring Workflow
Step-by-step workflow for creating an effective Agent Skill.

This workflow follows the guidance from [Agent Skills best practices](https://agentskills.io/skill-creation/best-practices), [quickstart](https://agentskills.io/skill-creation/quickstart), and [evaluating skills](https://agentskills.io/skill-creation/evaluating-skills). See also `standards://skill-authoring` for the rules that govern skill format and content.

## Source References

- **Quickstart**: https://agentskills.io/skill-creation/quickstart
- **Best Practices**: https://agentskills.io/skill-creation/best-practices
- **Optimizing Descriptions**: https://agentskills.io/skill-creation/optimizing-descriptions
- **Using Scripts**: https://agentskills.io/skill-creation/using-scripts
- **Evaluating Skills**: https://agentskills.io/skill-creation/evaluating-skills

---

# Step 1: Extract Real Expertise

Do not ask an LLM to generate a skill from scratch using only its general training knowledge — the result will be vague and generic. Effective skills are grounded in real, domain-specific expertise.

## Option A: Extract from a hands-on task

1. Complete a real task in conversation with an agent, providing context, corrections, and preferences along the way.
2. Extract the reusable pattern into a skill. Pay attention to:
   - **Steps that worked** — the sequence of actions that led to success
   - **Corrections you made** — places where you steered the agent's approach (e.g., "use library X instead of Y," "check for edge case Z")
   - **Input/output formats** — what the data looked like going in and coming out
   - **Context you provided** — project-specific facts, conventions, or constraints the agent didn't already know

## Option B: Synthesize from existing project artifacts

1. Gather project-specific source material: internal documentation, runbooks, style guides, API specifications, schemas, configuration files, code review comments, issue trackers, version control history (especially patches and fixes), and real-world failure cases and their resolutions.
2. Feed this material to an LLM and ask it to synthesize a skill.
3. Verify the output captures YOUR schemas, failure modes, and recovery procedures — not generic advice.

---

# Step 2: Draft the Skill

## Create the directory and `SKILL.md`

1. Create the skill directory. The directory name must follow the `name` field rules (lowercase, hyphens only, no consecutive hyphens).
2. Create `SKILL.md` with the required frontmatter fields: `name` and `description`.
3. Write the body content — the instructions the agent follows when the skill activates.

## Structure for progressive disclosure

1. Keep `SKILL.md` under 500 lines and under 5,000 tokens.
2. Move detailed reference material to `references/`, templates to `assets/`, and reusable scripts to `scripts/`.
3. Tell the agent WHEN to load each referenced file with conditional triggers (e.g., "Read references/api-errors.md if the API returns a non-200 status code").

## Write effective instructions

1. Add what the agent lacks; omit what it already knows. Cut content if the agent would handle it correctly without the instruction.
2. Give the agent freedom when multiple approaches are valid; be prescriptive when operations are fragile or a specific sequence must be followed.
3. Pick a default approach and mention alternatives briefly — never present all options as equal.
4. Teach the agent HOW to approach a class of problems, not WHAT to produce for a specific instance.
5. Include a **gotchas** section for environment-specific facts that defy reasonable assumptions.

## Include instruction patterns as needed

- **Checklists**: For multi-step workflows with dependencies or validation gates.
- **Output templates**: When the agent must produce output in a specific format.
- **Validation loops**: Do the work → run validation → fix issues → repeat until pass.
- **Plan-Validate-Execute**: For batch or destructive operations (create plan → validate → execute).

## Bundle scripts when needed

1. If the task involves repeated logic (building charts, parsing a format, validating output), write a tested script once and bundle it in `scripts/`.
2. Make scripts non-interactive — accept all input via CLI flags, environment variables, or stdin.
3. Include `--help` output, helpful error messages, and structured output (JSON/CSV/TSV to stdout, diagnostics to stderr).
4. Pin dependency versions for reproducibility.

---

# Step 3: Refine with Real Execution

1. Run the skill against real tasks.
2. Read agent execution traces — not just final outputs. Look for:
   - Instructions that are too vague (agent tries several approaches before finding one that works)
   - Instructions that don't apply to the current task (agent follows them anyway)
   - Too many options presented without a clear default
3. When an agent makes a mistake you have to correct, add the correction to the gotchas section.
4. Feed all results (not just failures) back into the creation process.
5. Repeat: even a single pass of execute-then-revise noticeably improves quality.

---

# Step 4: Optimize the Description

The description is the primary mechanism agents use to decide whether to load a skill. A poorly written description means the skill won't trigger when it should (or will trigger when it shouldn't).

## Write the initial description

1. Use imperative phrasing ("Use this skill when..." not "This skill does...").
2. Focus on user intent and what the user is trying to achieve.
3. Be explicit about when the skill applies, including cases where the user doesn't name the domain directly.
4. Include specific keywords.
5. Keep it under 1024 characters.

## Design eval queries

1. Create 8-10 should-trigger queries and 8-10 should-not-trigger queries in `evals/evals.json`.
2. Vary queries by phrasing, explicitness, detail, and complexity.
3. For should-not-trigger queries, prioritize **near-misses** — queries that share keywords with the skill but need something different.
4. Include realistic context: file paths, personal context, column names, casual language.

## Run the optimization loop

1. Split queries into train (~60%) and validation (~40%) sets.
2. Run each query multiple times (at least 3) to compute a trigger rate.
3. Use train set failures to guide description revisions (never use validation set for this).
4. Revise the description to generalize — find the underlying category of failing queries, not specific keywords.
5. Select the best iteration by its validation pass rate (not necessarily the last iteration).
6. Stop after ~5 iterations or when improvements plateau.

---

# Step 5: Evaluate Output Quality

Use structured evaluations (evals) to verify the skill produces good outputs reliably.

## Set up eval infrastructure

1. Create `evals/evals.json` in the skill directory with test cases (prompt, expected output, optional input files).
2. Start with 2-3 test cases; expand after seeing first results.
3. Vary prompts by phrasing, detail, and formality. Cover edge cases.
4. Include at least one boundary condition (malformed input, unusual request, ambiguous instructions).

## Run evals

1. For each test case, run twice: once WITH the skill and once WITHOUT it (baseline).
2. Use a clean context for each run — no leftover state.
3. Capture timing data (total tokens, duration) for each run.
4. Organize results in a workspace directory with `iteration-N/eval-name/with_skill/` and `without_skill/` subdirectories.

## Write assertions and grade

1. After seeing first-round outputs, add assertions to each test case — verifiable statements about what the output should contain.
2. Grade each assertion as PASS or FAIL with concrete evidence.
3. Aggregate results to compute pass rates, time, and token deltas.
4. Analyze patterns: remove assertions that always pass, investigate assertions that always fail, study assertions that pass with the skill but fail without, tighten instructions for inconsistent results.
5. Human review: check outputs alongside grades for quality issues assertions didn't catch.

## Iterate

1. Give failed assertions, human feedback, execution transcripts, and the current `SKILL.md` to an LLM and ask it to propose improvements.
2. Apply changes, rerun all test cases in a new `iteration-N/` directory, grade, review, and repeat.
3. Stop when satisfied, feedback is consistently empty, or no meaningful improvement between iterations.
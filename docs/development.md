# development.md

## Issue Rules

- 原則として1 Issue = 1つの明確な目的
- 1 Issueごとにfeature branchを作成する
- 原則として1 Issue = 1 PR
- Issueには以下を記載する
  - 目的
  - 現状
  - 実装内容
  - 対象外
  - 関連仕様
  - 完了条件
  - テスト項目
- 実装前に既存コードと影響範囲を調査する
- 実装後は関連するテスト・lint・buildなどのチェックを実行する
- 不要なリファクタリングは同じIssueに含めない
- Issueの実装前に実装方針を整理し、変更予定ファイルを確認する

## Claude Code workflow

1. ユーザーは Claude Code に粗い要求や目的を提示する。
2. Claude Code は `CLAUDE.md`、`docs/`、既存コード、既存 Issue を確認する。
3. Claude Code は以下を含む GitHub Issue 案を作成する。
   - 目的
   - 現状
   - 実装内容
   - 対象外
   - 関連仕様
   - 完了条件
   - テスト項目
4. ユーザーが Issue 案を確認し、GitHub Issue を作成する。
5. ユーザーが feature branch を作成・切り替える。
6. Claude Code は対象 Issue を読み、実装前調査を行う。
7. Claude Code は以下を整理してユーザーに提示する。
   - Issue の目的
   - 実装スコープ
   - 対象外
   - 変更予定ファイル
   - 影響範囲
   - 実装方針
   - テスト方針
   - 現在の Issue を進めるために必要な確認事項
8. ユーザーの確認・承認後に Claude Code が実装を開始する。
9. Claude Code は実装後に関連するテスト・チェックを実行する。
10. Claude Code は変更内容とテスト結果をまとめ、Git/GitHub 操作を行わずに停止する。
11. ユーザーが差分を確認し、commit / push / Pull Request 作成などを行う。

When information is missing, identify only the information required for the current issue.
Do not block implementation on unrelated future decisions.

## Git and GitHub responsibility

Claude Code is responsible for investigating, proposing, implementing, and testing changes.

The user is responsible for all Git and GitHub write operations.

### Claude Code may

The following read-only operations may be performed as needed:

- `git status`
- `git diff`
- `git log`
- `git branch`
- `git show`
- `git remote -v`

Claude Code may also:

- Investigate existing code and dependencies
- Inspect GitHub Issues
- Propose GitHub Issue content
- Propose implementation approaches
- Modify source files after user approval
- Run tests, lint, build, and other relevant checks
- Summarize changed files and diffs
- Suggest commit messages
- Suggest Pull Request titles and descriptions

### Claude Code must not

Claude Code must not perform any Git or GitHub write operation.

This includes:

- Creating or editing GitHub Issues
- Closing GitHub Issues
- Creating or switching branches
- Staging files with `git add`
- Creating commits
- Amending commits
- Pushing to remote repositories
- Creating, editing, merging, or closing Pull Requests
- Merging branches
- Deleting branches
- Rebasing
- Resetting commits
- Reverting commits
- Cherry-picking commits
- Force pushing
- Modifying Git history
- Performing any other operation that changes the local Git state, remote repository, or GitHub state

All Git and GitHub write operations are performed by the user.

## Pull Request review rules

Pull Requests use the template at `.github/PULL_REQUEST_TEMPLATE.md`.

Claude Code proposes Pull Request body content that follows this template.
Creating the Pull Request itself remains a user responsibility, as defined in "Git and GitHub responsibility" above.

When reviewing a Pull Request, check that:

- The change stays within the scope described in the linked Issue
- The Pull Request body accurately reflects the actual diff
- Required tests were run and their results are recorded
- Impact on other systems is described
- The points to review are clearly stated

## Destructive Git operations

Claude Code must never execute destructive Git operations.

Examples include:

- `git reset --hard`
- `git push --force`
- `git push --force-with-lease`
- `git clean -fd`
- `git branch -D`
- remote branch deletion

If such an operation appears necessary, Claude Code should explain:

- why it may be necessary
- what it would change
- what risks it has
- a safer alternative, if available

The user decides whether to execute the operation manually.

## After implementation

After implementation and testing, Claude Code should stop and report:

- changed files
- implementation summary
- test results
- relevant diff summary
- any remaining concerns or follow-up items
- suggested commit message
- suggested Pull Request title
- suggested Pull Request description

Claude Code must not stage, commit, push, create a Pull Request, or perform any other Git/GitHub write operation after reporting.
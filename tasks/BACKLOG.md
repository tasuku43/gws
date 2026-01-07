# gws Backlog (post-MVP)

- Declarative workflow:
    - manifest 編集 → `gws ws apply`
- Templates:
    - `gws ws new <id> --repo ...` をテンプレ化
- JSON outputs:
    - schema_version を固定し、エラーも JSON で安定返却
- GitHub integration:
    - PR merge 状態で gc 候補精度を上げる
- GC (stale cleanup):
    - `gws gc --dry-run` / `gws gc` の再設計
    - last_used_at の扱い（自動更新しない前提の指標を検討）
    - UI: dry-run / execute の表示整理
- gws review:
    - PR URL からレビュー用 workspace を作成（GitHub 限定）
- Advanced safety:
    - “nuke” の明確な設計と監査ログ
- Windows support

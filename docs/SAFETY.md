# 安全設計（破壊的操作のガード）

## 原則
- Safe by default
- dirty（未コミット変更）や未追跡ファイルがある場合、削除・回収を拒否する

## dirty 判定（MVP）
- `git status --porcelain` が非空なら dirty=true
- 未追跡ファイルを含めるかは設定で将来切替可能（MVPは含める）

## 破壊的フラグ
MVPでは破壊的フラグは提供しない

## 推奨運用
（現時点で該当なし）
- pinned workspace は回収しない

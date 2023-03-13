# vmmgr Controller
### How to start?
```
git submodule update --init --recursive
docker compose build
docker compose up -d
```

### 注意点
- qcow2のファイルは、rawファイルに変換する必要あり。
  - qcow2からrawに変換する実装に難ありのため
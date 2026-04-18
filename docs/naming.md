# Java → Go 名前対応表

書籍のJavaサンプル実装からGoで再実装する際に変更したクラス・メソッド名の対応表。

## file パッケージ

### 型

| Java クラス | Go の型 | 種別 | 備考 |
|---|---|---|---|
| `FileMgr` | `file.BlockStore` | interface | 使う側で定義するためinterfaceとして切り出し |
| `FileMgr` | `file.blockStore` | struct | `BlockStore` の実装 |
| `Page` | `file.Page` | interface | 使う側で定義するためinterfaceとして切り出し |
| `Page` | `file.page` | struct | `Page` の実装 |
| `BlockId` | `file.BlockID` | struct | Goの命名規則に合わせID→IDに変更 |

### メソッド・関数

| Java | Go | 備考 |
|---|---|---|
| `new FileMgr(dbPath, blockSize)` | `file.NewBlockStore(dbPath, blockSize)` | |
| `new Page(blockSize)` | `file.NewPage(blockSize)` | |
| `new Page(bytes)` | `file.NewPageFromBytes(bytes)` | バイト列からの生成を別関数に分離 |
| `new BlockId(fileName, blockNum)` | `file.NewBlockID(fileName, blockNum)` | |
| `FileMgr.length(fileName)` | `BlockStore.BlockCount(fileName)` | `length` は文字列長と混同しやすいため改名 |
| `Page.contents()` | `Page.Contents()` | Goの公開メソッドに合わせ大文字に変更 |
| `Page.maxLength(strlen)` | `file.MaxLength(strlen)` | インスタンスメソッド → パッケージレベル関数 |
| `Integer.BYTES` | `file.Int32Bytes` | Javaの定数をGoの定数として定義 |

## log パッケージ

### 型

| Java クラス | Go の型 | 種別 | 備考 |
|---|---|---|---|
| `LogMgr` | `log.Appender` | interface | 責務（追記）を名前に反映 |
| `LogMgr` | `log.appender` | struct | `Appender` の実装 |

### メソッド・関数

| Java | Go | 備考 |
|---|---|---|
| `new LogMgr(fm, logfile)` | `log.NewAppender(blockStore, logFile)` | |
| `LogMgr.append(logrec)` | `Appender.Append(logrec)` | |
| `LogMgr.flush(lsn)` | `Appender.Flush(lsn)` | |
| `LogMgr.iterator()` | `Appender.All()` | Javaのイテレータ → Go の `iter.Seq2` |

## buffer パッケージ

### 型

| Java クラス | Go の型 | 種別 | 備考 |
|---|---|---|---|
| `BufferMgr` | `buffer.Pool` | interface | 責務（バッファプール）を名前に反映 |
| `BufferMgr` | `buffer.pool` | struct | `Pool` の実装 |
| `Buffer` | `buffer.Buffer` | interface | |
| `Buffer` | `buffer.buffer` | struct | `Buffer` の実装 |

### メソッド・関数

| Java | Go | 備考 |
|---|---|---|
| `new BufferMgr(fm, lm, numBuffs)` | `buffer.NewPool(blockStore, appender, numBuffs)` | |
| `BufferMgr.pin(block)` | `Pool.Pin(block)` | Go版は `(Buffer, error)` を返す |
| `BufferMgr.unpin(buff)` | `Pool.UnPin(buff)` | |
| `BufferMgr.available()` | `Pool.Available()` | |
| `BufferMgr.flushAll(txnum)` | `Pool.FlushAll(txnum)` | |
| `new Buffer(fm, lm)` | `buffer.NewBuffer(blockStore, appender)` | |
| `Buffer.contents()` | `Buffer.Contents()` | |
| `Buffer.block()` | `Buffer.Block()` | |
| `Buffer.setModified(txnum, lsn)` | `Buffer.SetModified(txnum, lsn)` | |
| `Buffer.isPinned()` | `Buffer.IsPinned()` | |
| `Buffer.modifyingTx()` | `Buffer.ModifyingTx()` | |

## 作业

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## 解答
dao层应该Wrap这个error，抛给上层，由api层处理这个error，写入日志、生成对应的业务code返回给用户。
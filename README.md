# jump

jump 是一款本地管理 ssh 连接的工具，它的交互灵感来源于 JumpServer，区别在于 jump 用于管理本地要连接的主机，并非一款堡垒机工具。

所有的主机信息保存在 sqlite 中，路径位于 `~/.jump/host.db`。

该工具暂时属于自用，sqlite 中信息暂未实现加密存储。

## 截图
![](docs/img/jump-demo.jpg)

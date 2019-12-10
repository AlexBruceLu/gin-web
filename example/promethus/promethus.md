

# promethus监控系统

## 学习目标

- [ ] **能够安装promethus服务器**
- [ ] **能够安装node_exporter远程监控linux**
- [ ] **能够安装mysqld_exporter远程监控mysqld**
- [ ] **<font color=red>能够安装Grafana</font>**
- [ ] **<font color=red>能够在Grafana添加promethus数据源</font>**
- [ ] **<font color=red>能够在Grafana添加CPU负载的图形</font>**
- [ ] **<font color=red>能够在Grafana图形显示mysql监控数据</font>**
- [ ] **<font color=red>能够通过Grafana+onealert实现报警</font>**

## 任务背景

公司网站快速拓展，对现有项目进行监控。

## 任务要求

1. 部署监控服务器，实现7*24实时监控
2. 设计监控系统，对监控器和触发器有合理意见
3. 有好的预警机制，对可能出现的问题要及时警告并形成严格的处理机制
4. 做好监控系统，可以实现警告分级
   - 一级警报 电话通知
   - 二级警报 微信通知
   - 三级警报 邮件通知
5. 处理异地集中监控


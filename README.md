# PEAEDGE

PeaEdge是一个物联网边缘网关软件，可以运行在各种支持Linux的嵌入式设备上，支持 linux X86_64, ARM, ARM64, Mips。

## 系统架构

<img width="921" alt="image" src="https://user-images.githubusercontent.com/377938/198872775-85a16ce9-a36e-4ee7-b6d2-94f3319f259c.png">

## 功能特性

- [x] 支持标准Modbus设备访问管理
- [x] 实现 modbus tcp slave 和 rtu slave
- [x] 多数据通道支持， MQTT， TCP， HTTP 多协议协议支持
- [x] 控制流支持，通过下行数据触发控制流，实现对设备的控制操作
- [x] Hj212协议支持
- [x] LUA 脚本嵌入支持， 灵活自定义功能
- [x] 告警支持，支持 LUA 脚本自定义告警规则
- [x] 采用 SQLite 数据库存储数据，小巧灵活 资源占用少

### [- Wiki 文档手册](https://github.com/peaiot/peaedge/wiki)

## 许可证

参见 [GNU LESSER GENERAL PUBLIC LICENSE](LICENSE)

- [安装说明]()

## 贡献

我们欢迎任何形式的贡献，包括但不限于问题、拉动请求、文档、例子等。

## 联系我们

如果你有任何问题，你可以通过电子邮件联系我们

# QuickClash

QuickClash是一个用于快速启动Clash的工具，适用于临时设备，并且仅支持Windows操作系统。

## 关于mmdb问题

如果以前没用过clash，那么启动会失败，因为无mmdb且由于网络原因下载失败。

可以下载[补充脚本](https://github.com/UlinoyaPed/QuickClash/releases/download/v1.3/mmdb.zip)以解决，内含2024.5.11的IP数据库（非最新的，但至少能运行起来）。

## 功能特点

### 程序流程：
1. 下载配置文件和WebUI资源文件。
2. 加载和解析配置文件。
3. 启动一个Gin实例，用于提供Web界面。
4. 检查订阅文件是否需要更新，并在需要时更新订阅。
5. 检查所需的Clash内核文件是否存在，并在需要时下载。
6. 启动Clash，并使用配置文件运行。
7. 设置系统代理，将流量导向Clash。
8. 在程序终止时自动取消代理。
9. 打印相关提示信息。

### 特点：
1. 支持多语言功能。
2. 通过配置文件选择使用的Clash内核（Clash或ClashMeta）。
3. 支持下载和解压WebUI资源文件。
4. 支持自动检查和下载新的配置文件。
5. 在启动Clash之后设置系统代理。
6. 使用Gin作为Web界面的服务器。

### NEW:
- 进度条
- 支持32位系统

## 使用方法

1. 下载最新版本的QuickClash工具。
2. 解压缩下载的文件到你想要存放的目录。
3. 双击运行QuickClash.exe文件。
4. Clash将会启动。

## 注意事项

- 本工具仅适用于Windows操作系统。

## 许可证

本项目使用GPL-3.0许可证。有关详细信息，请参阅 [LICENSE](LICENSE) 文件。

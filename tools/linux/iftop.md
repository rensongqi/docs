
# iftop 高级操作

1. 启动时直接过滤IP

    在运行 iftop 时，通过 -f 参数指定过滤条件：
    ```bash
    sudo iftop -f 'host 192.168.1.100'      # 过滤特定 IP
    sudo iftop -f 'host example.com'        # 过滤域名（需支持 DNS 解析）
    ```

2. 运行时交互式过滤

    在 iftop 运行界面中，通过快捷键动态过滤：

    > 按 l 键，输入过滤表达式（如 host 192.168.1.100），回车确认。
    > 
    > 按 L 键切换显示/隐藏流量柱状图。
    > 
    > 按 p 键切换显示端口信息。

3. 按源/目的地址过滤

    > 过滤源地址（Source）：按 s 键，输入源 IP 或域名。
    >  
    > 过滤目的地址（Destination）：按 d 键，输入目标 IP 或域名。
    > 
    > 同时过滤源和目的：按 l 键输入 src host <IP> and dst host <IP>。
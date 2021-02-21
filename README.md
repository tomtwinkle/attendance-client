# 奉行クラウド打刻 & Slack通知 CLI
Attendance Bugyo Cloud Punchmark & post Slack message CLI

# Installation

## Windows
- Install with PowerShell

```poweshell
iwr https://github.com/tomtwinkle/attendance-client/releases/download/v0.0.1/attendance_windows_amd64.zip -OutFile attendance.zip && Expand-Archive -Path attendance.zip && rm attendance.zip
cd attendance
.\attendance.exe help
```

# How to use CLI

- 出勤

```shell
.\attendance punchmark --type in
or 
.\attendance pm -t in
```

- 退出

```shell
.\attendance punchmark --type out
or 
.\attendance pm -t out
```

- 外出

```shell
.\attendance punchmark --type go
or 
.\attendance pm -t go
```

- 再入

```shell
.\attendance punchmark --type return
or 
.\attendance pm -t return
```

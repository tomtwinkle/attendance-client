# 奉行クラウド打刻 & Slack通知 CLI
Attendance Bugyo Cloud Punchmark & post Slack message CLI

# Installation

## Windows
- Install with PowerShell

```poweshell
iwr https://github.com/tomtwinkle/attendance-client/releases/download/v0.0.6/attendance_windows_amd64.zip -OutFile attendance.zip && Expand-Archive -Path attendance.zip && rm attendance.zip
cd attendance
.\attendance.exe help
```

- Generate Slack Token

> 1. https://api.slack.com/apps を開く
> 2. [Create New App] を選択
> 3. [OAuth & Permissions] を選択
> 4. [User Token Scopes] を選択 [channels:read, groups:read, chat:write] を追加
> 5. [Install to Workspace] を選択
> 6. OAuth Access Token [Copy] を選択

- Add App to Slack channel

チャンネルにメッセージを投稿するためにはチャンネルにアプリを追加する必要があります.

> 1. [(i) チャンネル詳細] を選択
> 2. [... その他] を選択
> 3. [アプリを追加する] を選択
> 4. 作成したアプリを検索し、[追加]を選択. [表示する]となっている場合は追加は不要です.

# How to use CLI

## 設定ファイル作成

- 初期設定

```shell
.\attendance init
```

- 作成された `attendance.yaml` を任意で編集

```
tenantcode: <テナントコード*1>
obcid: <OBC ID>
password: <OBC Password>
slack: *2
    token: <slack token>
    channels:
      - name: <slackチャンネル名*3>
        clockin:
            message: "<出勤時に投稿するコメント*4>"
        clockout:
            message: "<退勤時に投稿するコメント*4>"
        goout:
            message: "<外出時に投稿するコメント*4>"
        returned:
            message: "<再入時に投稿するコメント*4>"
```

*1 : URLを参照. `https://id.obc.jp/<テナントコード>`

*2 : slack 要素の削除でslackへの投稿は行われなくなります.

*3 : yaml Sequence記法で複数チャンネルへ投稿することも可能です.
[Sequence of Mappingsを参照](https://yaml.org/spec/1.2/spec.html#id2759963)

*4 : コメントをbankにする or clockin,clockout等の要素削除でslackへの投稿は行われなくなります.

## 出勤

```shell
.\attendance punchmark --type in
or 
.\attendance pm -t in
```

## 退出

```shell
.\attendance punchmark --type out
or 
.\attendance pm -t out
```

## 外出

```shell
.\attendance punchmark --type go
or 
.\attendance pm -t go
```

## 再入

```shell
.\attendance punchmark --type return
or 
.\attendance pm -t return
```

## その他オプション

- slackにのみ投稿する

```shell
.\attendance punchmark --type in --slackonly
```

- 打刻のみ行う

```shell
.\attendance punchmark --type in --slackskip
```

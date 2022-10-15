---
url: /en/usage/
title: Usage examples
weight: 10
---

After you [install the CLI][install] the first thing to do is run the command `clockifycli config
init`, it will interactively ask you the information necessary to setup your environment.

```console
$ clockify-cli config init
? User Generated Token: <your-api-token>
? Choose default Workspace: <workspace-id> - John Doe's workspace
? Choose your user: <user-id> - John Doe
? Should try to find projects/tasks/tags by their names? Yes
? Should use "Interactive Mode" by default? Yes
? Which days of the week do you work? monday, tuesday, wednesday, thursday, friday
? Should allow starting time entries with incomplete data? No
? Should show task on time entries as a separated column? Yes
```

These answers will be saved at `$HOME/.clockify-cli.yaml` by default and can be copied from machine
to machine if needed. You can see all the options [here][cli-config].

> ❗ If you have installed the client using `snap` this file will not be accessible to
> you, but the configurations will still be persisted.

After that you can start a new entry using the `clockify-cli in` command.

```console
$ clockify-cli in
? Choose your project: 621948458cb9606d934ebb1c - Clockify Cli      | Client: Myself (6202634a28782767054eec26)
? Choose your task: 62ae29e62518aa18da2acd14 - In Command
? Description: Some description
? Choose your tags: 62ae28b72518aa18da2acb49 - Development
? Start: 2022-06-30 22:49:34
? End (leave it blank for empty):
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
|            ID            |  START   |   END    |   DUR   |   PROJECT    |   DESCRIPTION    |                  TAGS                  |
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
| 62be52d2f2c0e80ba36fce0a | 22:49:34 | 22:50:10 | 0:00:36 | Clockify Cli | Some description | Development (62ae28b72518aa18da2acb49) |
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
```

> See more about [in here][cli-in]

By default it will prompt you about the details of the time entry, if you don't like that as the
default behavior you can change it by running the command `clockify-cli config interactive false`,
and if there is a situation were you want to be prompted, then run `clockify-cli in --interactive`.

This behavior is true for all interactive commands, except for `clockify-cli config init`.

Once you finish the activity or need to stop timer, run `clockify-cli out` to stop it.

```console
$ clockify-cli out
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
|            ID            |  START   |   END    |   DUR   |   PROJECT    |   DESCRIPTION    |                  TAGS                  |
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
| 62be52d2f2c0e80ba36fce0a | 22:49:34 | 23:02:41 | 0:13:07 | Clockify Cli | Some description | Development (62ae28b72518aa18da2acb49) |
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
```
> See more about [out here][cli-out]

To start a new timer with the same information as the last one you did, you can run `clockify-cli
clone last` and a new timer with the same properties as the last stopped one will be started.

```console
$ clockify-cli clone last -i=0 # -i=0 will stop the CLI from prompting you about the timer
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
|            ID            |  START   |   END    |   DUR   |   PROJECT    |   DESCRIPTION    |                  TAGS                  |
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
| 62be5684f2c0e80ba36fcecd | 23:05:53 | 23:05:56 | 0:00:03 | Clockify Cli | Some description | Development (62ae28b72518aa18da2acb49) |
+--------------------------+----------+----------+---------+--------------+------------------+----------------------------------------+
```

Lets say that the current activity is for the same task, project and description, but you doing a
pairing with someone now. You can fix the timer using the `clockify-cli edit current` command and
change the tags of the timer.

```console
$ clockify-cli edit current -T pair -T web
+--------------------------+----------+----------+---------+--------------+------------------+---------------------------------------------+
|            ID            |  START   |   END    |   DUR   |   PROJECT    |   DESCRIPTION    |                    TAGS                     |
+--------------------------+----------+----------+---------+--------------+------------------+---------------------------------------------+
| 62be5684f2c0e80ba36fcecd | 23:05:53 | 23:29:14 | 0:23:21 | Clockify Cli | Some description | Pair Programming (621948708cb9606d934ebba7) |
|                          |          |          |         |              |                  | Development (62ae28b72518aa18da2acb49)      |
+--------------------------+----------+----------+---------+--------------+------------------+---------------------------------------------+
```
> See more about [edit here][cli-edit]

Now you remembered that yesterday you had a meeting at the end of the day that you forgot to
register, but you don't want to stop the running one.

To create a time entry that has a start and end without tempering with a running timer you can use
the command `clockify-cli manual`.

```console
$ clockify-cli manual -s "yesterday 17:50" -e "yesterday 18:00" -T meet -d 'About the Calendar' \
    -p cli

+--------------------------+----------+----------+---------+--------------+--------------------+------------------------------------+
|            ID            |  START   |   END    |   DUR   |   PROJECT    |    DESCRIPTION     |                TAGS                |
+--------------------------+----------+----------+---------+--------------+--------------------+------------------------------------+
| 62be5d49f2c0e80ba36fd01e | 17:50:00 | 18:00:00 | 0:10:00 | Clockify Cli | About the Calendar | Meeting (6219485e8cb9606d934ebb5f) |
+--------------------------+----------+----------+---------+--------------+--------------------+------------------------------------+
```
> See more about [manual here][cli-manual]

If you forgot to stop a timer running and wants to stop it with a specific instead of now, you can
use the flag `--when` to set the end time.

```console
$ clockify-cli out --when 23:35
+--------------------------+----------+----------+---------+--------------+------------------+---------------------------------------------+
|            ID            |  START   |   END    |   DUR   |   PROJECT    |   DESCRIPTION    |                    TAGS                     |
+--------------------------+----------+----------+---------+--------------+------------------+---------------------------------------------+
| 62be5684f2c0e80ba36fcecd | 23:05:53 | 23:35:00 | 0:29:07 | Clockify Cli | Some description | Pair Programming (621948708cb9606d934ebba7) |
|                          |          |          |         |              |                  | Development (62ae28b72518aa18da2acb49)      |
+--------------------------+----------+----------+---------+--------------+------------------+---------------------------------------------+
```

To see the entries for today you can use the command `clockify-cli report` to list them.

```console
$ clockify-cli report
+--------------------------+---------------------+---------------------+---------+--------------+-------------------+---------------------------------------------+
|            ID            |        START        |         END         |   DUR   |   PROJECT    |    DESCRIPTION    |                    TAGS                     |
+--------------------------+---------------------+---------------------+---------+--------------+-------------------+---------------------------------------------+
| 62be52d2f2c0e80ba36fce0a | 2022-06-30 22:49:34 | 2022-06-30 23:02:41 | 0:13:07 | Clockify Cli | Some description  | Development (62ae28b72518aa18da2acb49)      |
+--------------------------+---------------------+---------------------+---------+--------------+-------------------+---------------------------------------------+
| 62be5684f2c0e80ba36fcecd | 2022-06-30 23:05:53 | 2022-06-30 23:35:00 | 0:29:07 | Clockify Cli | Some description  | Pair Programming (621948708cb9606d934ebba7) |
|                          |                     |                     |         |              |                   | Development (62ae28b72518aa18da2acb49)      |
+--------------------------+---------------------+---------------------+---------+--------------+-------------------+---------------------------------------------+
| 62be5eb535710e76ef03c884 | 2022-06-30 23:40:42 | 2022-06-30 23:41:27 | 0:00:45 | Clockify Cli | Other description | Development (62ae28b72518aa18da2acb49)      |
+--------------------------+---------------------+---------------------+---------+--------------+-------------------+---------------------------------------------+
| TOTAL                    |                     |                     | 0:42:59 |              |                   |                                             |
+--------------------------+---------------------+---------------------+---------+--------------+-------------------+---------------------------------------------+
```
> See more about [report here][cli-report]

If you need to quickly see how much time was spent this month in a project you can use
`clockify-cli report this-month`, the flag `--project` to filter the timers and
`--duration-formatted` to get only the sum of time.

```console
$ clockify-cli report this-month -p cli --duration-formatted
6:23:52
```

[install]: https://github.com/lucassabreu/clockify-cli#how-to-install-
[cli-config]: /en/commands/clockify-cli_config/
[cli-in]: /en/commands/clockify-cli_in/
[cli-manual]: /en/commands/clockify-cli_manual/
[cli-clone]: /en/commands/clockify-cli_clone/
[cli-report]: /en/commands/clockify-cli_report/
[cli-edit]: /en/commands/clockify-cli_edit/

# Go daemon

This repo is designated to create a daemon in go following the article:
[Four steps to Daemonize your Go Programs](https://ieftimov.com/posts/four-steps-daemonize-your-golang-programs/)
Daemon's installer was created following this other article: [Create and Manage Launchd Agents Macos](https://ieftimov.com/posts/create-manage-macos-launchd-agents-golang/)

## Four steps are
- 1 Log to standard output
- 2 Shutdonw on SIGTERM/SIGINT
- 3 Reload the config on SIGHUP
- 4 Provide necessary config file for your favorite init system to control your `daemon`

In the code you can find comments that indicates which the parts that are implementing those four steps.



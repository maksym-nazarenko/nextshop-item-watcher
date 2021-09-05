Next online shop items watcher
==============================

An educational project to learn more about [Go](https://golang.org/) programming language.

This projects aims to help you with items on the [Next website](https://www.next.co.uk) that are in "ComingSoon", "SoldOut" and other statuses and send notification once the item is "InStock".

Currently, only [Telegram](https://telegram.org/) is supported as interface for adding watch items and receiving notifications.


## Development environment
1. Clone this repo
1. [Desireable] Link commit hooks from `git-hooks/` folder into your local `.git/hooks/` folder:

```shell
$ pushd .git/hooks/ && ln -s ../../git-hooks/* ./ && popd
```

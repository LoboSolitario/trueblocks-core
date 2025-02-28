# Recent Migrations

Please complete any applicable [Migrations](https://github.com/TrueBlocks/trueblocks-core/blob/develop/MIGRATIONS.md).

# TrueBlocks Core

![GitHub repo size](https://img.shields.io/github/repo-size/TrueBlocks/trueblocks-core)
[![GitHub contributors](https://img.shields.io/github/contributors/TrueBlocks/trueblocks-core)](https://github.com/TrueBlocks/trueblocks-core/contributors)
[![GitHub stars](https://img.shields.io/github/stars/TrueBlocks/trueblocks-core?style%3Dsocial)](https://github.com/TrueBlocks/trueblocks-core/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/TrueBlocks/trueblocks-core?style=social)](https://github.com/TrueBlocks/trueblocks-core/network/members)
[![Twitter Follow](https://img.shields.io/twitter/follow/trueblocks?style=social)](https://twitter.com/trueblocks)

## Table of Contents
  - [Introduction](#introduction)
  - [Quick Install](#quick-install)
  - [Introducing Chifra](#introducing-chifra)
  - [Using Chifra](#using-chifra)
  - [Building the Index](#building-the-index)
  - [Docker](#docker)
  - [Contributing](#contributing)
  - [List of Contributors](#list-of-contributors)
  - [Contact](#contact)

## Introduction

TrueBlocks is a collection of libraries, tools, and applications that improve access to Ethereum data while remaining entirely local. The interface is similar to the Ethereum RPC but offers several improvements:

1) TrueBlocks allows you to scrape the chain to build an index of address appearances. This index enables lightning-fast access to transactional histories for a given address (something not available from the node itself),

2) TrueBlocks also provides for a local binary cache of data extracted from the node. This speeds up subsequent queries for the same data by  order of magnitude or more. This enables a much better user experience for distributed applications written directly against the node, such as the [TrueBlocks Explorer](https://github.com/TrueBlocks/trueblocks-explorer),

3) TrueBlocks enhances the Ethereum RPC interfaces. For example, you may query blocks and transactions by date, by block range, by hashes, or any combination. Furthermore, two additional endpoints are provided to extract (`export`) and list (`list`) historical transactions per address.

## Quick install

TrueBlocks runs on Linux and Mac. There is no official Windows support.
Some users have had success using WSL─you're on your own!

These instructions assume you can navigate directories with the command line
and edit configuration files.
If you need help with a step, see the [installation's troubleshooting section](https://trueblocks.io/docs/install/install-trueblocks/#troubleshooting).

0. Install dependencies
    - &#9745; [Install the latest version of Go](https://golang.org/doc/install).
    - &#9745; Install the other dependencies with your command line: `build-essential` `git` `cmake` `ninja` `python` `python-dev` `libcurl3-dev` `clang-format` `jq`.

Alternatively, for nix users, you can drop into an isolated environment with necessary dependencies with `nix-shell`.

1. Compile from the codebase
    ```shell
    git clone -b develop https://github.com/trueblocks/trueblocks-core
    cd trueblocks-core
    mkdir build && cd build
    cmake ../src
    make
    ```
    _(You may use `make -j <ncores>` to parallelize the build. <ncores> represents the number of cores to devote to the `make` job)_

2. Add `trueblocks-core/bin` to your shell `PATH`.

3. Find your TrueBlocks configuration directory. It should be in one of these places:

    * On linux at `~/.local/share/trueblocks`
    * On mac at `~/Library/Application Support/TrueBlocks`
    * If you've configured it, wherever the location of `$XDG_CONFIG_HOME` is

4. In the configuration directory, edit `trueblocks.toml` to add your RPC and API keys. It should look something like this:
```toml
[settings]
rpcProvider = "<url-to-rpc-endpoint>"
```

5. Test a command!
```shell
chifra blocks 12345
```
### Optional steps

6. To make deep data queries, [get the index](https://trueblocks.io/docs/install/get-the-index/)
7. To explore the data visually, [install the explorer application](https://trueblocks.io/docs/install/install-explorer/).

## Introducing chifra

Similar to `git`, TrueBlocks has an overarching command called `chifra` that gives you access to all of the other subcommands.

Type:

```shell
chifra
```

You will see a long list of commands.

You can get more help on any command with `chifra <cmd> --help`.

### Getting status

Let's look at the first subcommand, called `status`.

```shell
chifra status --terse
```

If you get a bunch of JSON data, congratulations, your installation is working. You may skip ahead to the 'Using TrueBlocks' section below.

### -- Troubleshooting

Depending on your setup, you may get the following error message when you run some `chifra` commands:

```shell
  Warning: A request to your Ethereum node (http://localhost:8545) resulted
  in the following error [Could not connect to server]. Specify a valid
  rpcProvider by editing $CONFIG/trueblocks.toml.
```

If you get this error, edit the configuration file mentioned. The file is well documented, so refer to that file for further information.

When the `chifra status` command returns a valid response, you may move to the next section. If
you continue to have trouble, join our [discord discussion](https://discord.gg/kAFcZH2x7K).

## Using Chifra

If you've gotten this far, you're ready to use TrueBlocks. Run this command which shows every 10th block between the first and the 100,000th.

```shell
chifra blocks 0-100000:10
```

Hit `Control+C` to stop the processing.

## Building the index

---

The primary data structure produced by TrueBlocks is an index of address appearances. This index provides very quick access to transaction histories for a given address.

You may either build the entire index from scratch (requires a tracing, archive node) or download part of the index and build from there.

This process is described in this article [Indexing Addresses](https://trueblocks.io/docs/install/get-the-index/).

## Docker

Our official docker version is in a [separate repo](https://github.com/TrueBlocks/trueblocks-docker). Please see that repo for more information.

## Documentation

See our [documentation repo](https://github.com/TrueBlocks/trueblocks-docs) for the TrueBlocks website.


## Stargazers over time

A chart showing the number of stars on our repo over time.

[![Stargazers over time](https://starchart.cc/TrueBlocks/trueblocks-core.svg)](https://starchart.cc/TrueBlocks/trueblocks-core)

## Contributing

We love contributors. Please see information about our [work flow](https://github.com/TrueBlocks/trueblocks-core/blob/develop/docs/BRANCHING.md) before proceeding.

1. Fork this repository into your own repo.
2. Create a branch: `git checkout -b <branch_name>`.
3. Make changes to your local branch and commit them to your forked repo: `git commit -m '<commit_message>'`
4. Push back to the original branch: `git push origin TrueBlocks/trueblocks-core`
5. Create the pull request.

## List of Contributors

Thanks to the following people who have contributed to this project:

* [@tjayrush](https://github.com/tjayrush)
* [@dszlachta](https://github.com/dszlachta)
* [@wildmolasses](https://github.com/wildmolasses)
* [@MysticRyuujin](https://github.com/MysticRyuujin)
* [@MattDodsonEnglish](https://github.com/MattDodsonEnglish)
* [@crodnun](https://github.com/crodnun)

## Contact

If you have questions, comments, or complaints, please join the discussion on our discord server which is [linked from our website](https://trueblocks.io).

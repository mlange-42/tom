# Tom -- Terminal for Open-Meteo

An interactive terminal user interface (TUI) weather app.
Based on data from [Open-Meteo](https://open-meteo.com/).

Very early work in progress!

![screenshot forecast](https://github.com/mlange-42/tom/assets/44003176/83de75e0-babe-4a7a-8002-0ad4e2851855)
![screenshot plots](https://github.com/mlange-42/tom/assets/44003176/b0cc58d8-0565-43ae-ac81-2a2300e1e9bf)

## Installation

Pre-compiled binaries for Linux, Windows and MacOS are available in the
[Releases](https://github.com/mlange-42/tom/releases).

> Alternatively, install the latest development version of Tom using [Go](https://go.dev):
> ```shell
> go install github.com/mlange-42/tom@main
> ```

## Usage

Run `tom` with the name of a place:

```
tom Buxtehude
```

Set the default location, so you can later simply run `tom` without arguments:

```
tom Buxtehude --default
```

Get help:

```
tom -h
```

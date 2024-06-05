# Tom -- Terminal for Open-Meteo

An interactive terminal user interface (TUI) weather app.
Based on data from [Open-Meteo](https://open-meteo.com/).

Very early work in progress!

![screenshot forecast](https://github.com/mlange-42/tom/assets/44003176/d04c111c-9c20-4e95-b484-e487321e8578)
![screenshot plots](https://github.com/mlange-42/tom/assets/44003176/2db3b5cc-5256-4c7e-874f-39ed5ff7e5f0)

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

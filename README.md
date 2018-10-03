# ledger-go

ledger-go is a command-line application to work with double-entry accounting data stored in plain text.

## Install

```sh
go install github.com/kolo/ledger-go
```

## Data structure

ledger-go reads data from a folder specified by `LEDGER_DIR` environment variable. This folder has a specific structure. Top-level folders are named after an accounting year. A year folder has up to 12 month folders. A month folder contains text files named after days.

The data directory also contains a configuration file named `config.json`.

Here is an example layout of data directory:

```sh
./
.//
2018/
2018/01/
2018/02/
...
2018/10/
2018/10/01
2018/10/02
...
config.json
```

Data files are CSV files where each line presents a record. Each record has 3 values:

* credit account
* debit account
* amount

Here is an example of data file:

```txt
Card,Wallet,120
Card,Cafe,25.4
Wallet,Shop,12.75
```

## Configuration

Configuration is stored in a `config.json` file which is located in `LEDGER_DIR` directory. This is a JSON file which contains following settings:

* `assets` - list of asset accounts.

## Usage

Use `ledger-go --help` to get full list of supported commands.

## Contribution

Feel free to fork the project, submit pull requests, ask questions.

## Authors

Dmitry Maksimov (dmtmax@gmail.com)
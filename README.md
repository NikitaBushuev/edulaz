# EduLaZ Cryptocurrency

Welcome to the EduLaZ Cryptocurrency repository! This repository contains the Go implementation of EduLaZ Cryptocurrency, a simple cryptocurrency project designed for educational purposes. The project includes two executables: `edulaz-d` and `edulaz-cli`.

## Table of Contents

- [EduLaZ Cryptocurrency](#edulaz-cryptocurrency)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Executables](#executables)
    - [edulaz-cli Commands](#edulaz-cli-commands)
  - [Getting Started](#getting-started)
  - [Contribution](#contribution)
  - [License](#license)

## Introduction

EduLaZ Cryptocurrency is a learning-oriented project aimed at providing an introduction to the concepts of blockchain and cryptocurrencies. This implementation is written in Go and serves as a practical example to help you understand the underlying principles of blockchain technology.

## Executables

The repository contains two executable files that you can use to interact with the EduLaZ Cryptocurrency system:

1. `edulaz-d`: This executable is used to set up and run a node of the EduLaZ Cryptocurrency network. It requires two arguments: the `<path>` directory where it will open or create necessary files, and the `<address>` to identify the node. It creates `blockchain.json` and `private_key.json` files in the specified directory.

2. `edulaz-cli`: This executable provides a command-line interface to interact with the EduLaZ Cryptocurrency network. It requires the `<path>` argument, which should point to the directory containing the necessary blockchain and private key files.

The repository also includes a `Makefile` that simplifies common tasks. You can use the provided commands to streamline the process of working with the project.

### edulaz-cli Commands

The `edulaz-cli` executable supports the following commands:

- `mybalance`: Displays the balance of the current user's address.
- `myaddress`: Displays the address of the current user.
- `tx <address> <amount>`: Creates and sends a transaction of the specified amount to the given address.
- `balance <address>`: Displays the balance of the specified address.

Refer to the [Getting Started](#getting-started) section for information on how to use these commands.

## Getting Started

To get started with EduLaZ Cryptocurrency, follow these steps:

1. Clone this repository to your local machine.
2. Compile the executables using Make: `make all`.

3. Create a directory to store node files

    ```bash
    mkdir cli_9090 # node 1
    mkdir cli_9091 # node 2
    ```

4. Create the `addresses.json` file with a list of known IP node addresses.

    ```json
    // cli_9090/addresses.json
    [":9091"]

    // cli_9091/addresses.json
    [":9090"]
    ```

5. Run the `edulaz-d` executable to start a node, providing the necessary path and address arguments. The first running node (creator of the blockchain) receives 1024 `ELZ`

    ```bash
    edulaz-d ./cli_9090 :9090 # node 1 with ip :9090
    edulaz-d ./cli_9091 :9091 # node 2 with ip :9091
    ```

6. Use the `edulaz-cli` executable along with the available commands to interact with the network using the specified path.

    ```bash
    edulaz-cli ./cli_9090 # manage node 1 balance
    edulaz-cli ./cli_9091 # manage node 2 balance
    ```

Remember that this project is for educational purposes, and you can use it to learn about blockchain concepts and their implementation in Go.

## Contribution

Contributions to EduLaZ Cryptocurrency are welcome! If you find any issues or want to enhance the project, feel free to create pull requests or open issues in this repository.

## License

This project is licensed under the [MIT License](LICENSE).

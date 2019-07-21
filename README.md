# iot-abe

Attribute-based encryption and access control for the IoT.

## Contents <!-- omit in toc -->

- [iot-abe](#iot-abe)
  - [About](#about)
  - [Components](#components)
  - [Development](#development)
    - [Environment](#environment)
    - [Build](#build)
    - [Debugging](#debugging)
    - [Caveats](#caveats)
  - [Start a private Ethereum network](#start-a-private-ethereum-network)
  - [Raspberry Pi Integration](#raspberry-pi-integration)
  - [CLI](#cli)
    - [Examples](#examples)
  - [Links](#links)

## About

More and more electronic devices find their way into our everyday lives - for instance the coffee machine that automatically orders repair service, the parking lot that identifies when it is being used or the rubbish bin that detects when it is full. These devices are all connected to the Internet and are commonly denoted as the `Internet of Things (IoT)`. Most of the time they process and exchange data without any human intervention which introduces new concerns in terms of the security and privacy of the data.

The issues concerning the security and privacy of IoT systems turn out to be even more problematic in respect of the recent incidents that the IoT has [witnessed](https://www.iotforall.com/5-worst-iot-hacking-vulnerabilities).  
The main problem is that new technical challenges arise from the use of the IoT that are hard to resolve using traditional approaches. Some factors that contribute to this fact are:

- **heterogeneity**: The decentralized and highly heterogeneous nature of the IoT makes it difficult to centralize components such as authentication and authorization systems. This hints the use of a decentralized system such as Blockchain technologies to control both authentication and access control for the IoT.
- **credibility**: Not only the access control to the data but also the integrity of data sent by an IoT device has to be ensured.
- **harmonization**: IoT standards and protocols are only emerging slowly which might be due to the fact that the industry is developing faster than research in this area.
- **large volumes of data**: IoT devices are already producing a substantial amount of data which is expected to grow in the next few years. Approaches that are fully centralized will most probably soon reach the scaling limit and impose single points of failures. On the other hand, also approaches that are purely decentralized and try to save the data produced by the devices on a Blockchain (such as IOTA might soon reach their maximum capacity.
- **constrained resources**: IoT devices usually are resource constrained which prevents them from executing computationally intensive tasks.

To overcome the aforementioned issues as well as the fact that traditional access control models are hard to manage on scale. We propose using the rather new attribute-based access control model which leverages the fact that users are able to decrypt data based on a set of attributes. We introduce a new scheme for access control using several smart-contracts to facilitate ABE based on the Ethereum platform.

## Components

- Wrapper for the C library libbswabe
- Wallet management
- ABE Workflow
- IoT Server
- IoT Client

## Development

### Environment

Install the following components on your system:

- [PBC Crypto Library](https://crypto.stanford.edu/pbc/download.html)

  ```{.sh}
  sudo apt-get install flex bison
  cd output
  wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz -O ./pbc-0.5.14.tar.gz
  tar -zxvf pbc-0.5.14.tar.gz
  cd pbc-0.5.14/
  ./configure
  make
  sudo make install
  cd ..
  wget http://hms.isi.jhu.edu/acsc/cpabe/libbswabe-0.9.tar.gz -O ./libbswabe-0.9.tar.gz
  tar -zxvf libbswabe-0.9.tar.gz
  cd libbswabe-0.9
  ./configure
  make
  sudo make install
  cd ..
  wget http://hms.isi.jhu.edu/acsc/cpabe/cpabe-0.11.tar.gz -O ./cpabe-0.11.tar.gz
  tar -zxvf cpabe-0.11.tar.gz
  cd cpabe-0.11
  ./configure
  # add -lgmp to linker command, fix  policy_lang.y by changing line 67 to result: policy { final_policy = $1; }
  make
  sudo make install
  ```

To install all required development tools run:

```{.sh}
make setup
```

run:

```{.sh}
cp config.yml.dist config.yml
```

Create an [Etherscan api key](https://etherscan.io/apis) and add it to `config.yml`. Adapt the node urls in `config.yml` with the urls of your Ethereum nodes.

### Build

```{.sh}
make install
make build
```

### Debugging

go-delve does not support input from `stdin` when using vscode, therefore a remote debug session has to be launched.
Change the cmd flags in `.vscode/tasks.json` if needed.

```{.sh}
Ctrl + shift + B
Then press F5 (having Go Remote Debug selected)
```

Refer to:

- [Github Issue](https://github.com/Microsoft/vscode-go/issues/219)
- [Remote Debugging](https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code)

### Caveats

- `ecrecover` does not yield the correct result when using Ganache (see [Getting wrong address back](https://ethereum.stackexchange.com/questions/12621/getting-the-wrong-address-back-from-ecrecover/12684) and [ecrecover from web3](https://ethereum.stackexchange.com/questions/15364/ecrecover-from-geth-and-web3-eth-sign))

## Start a private Ethereum network

geth --datadir="/Users/Merryous/go/src/github.com/mariusgiger/iot-abe/contract/tmp/eth/" -verbosity 6 init genesis.json console 2>> ./tmp/eth/01.log

admin.nodeInfo.enode

geth --datadir="/Users/Merryous/go/src/github.com/mariusgiger/iot-abe/contract/tmp/eth/" -verbosity 6 --networkid 15 bootnode --genkey=boot.key console 2>> ./tmp/eth/01.log
geth --datadir="/Users/Merryous/go/src/github.com/mariusgiger/iot-abe/contract/tmp/eth/" -verbosity 6 --networkid 15 bootnode --nodekey=boot.key console 2>> ./tmp/eth/01.log

## Raspberry Pi Integration

- [Getting started](https://projects.raspberrypi.org/en/projects/raspberry-pi-getting-started/2)
- [Download Raspbian](https://www.raspberrypi.org/downloads/raspbian/)
- [Installing the Camera](https://thepihut.com/blogs/raspberry-pi-tutorials/16021420-how-to-install-use-the-raspberry-pi-camera)

Setup:

- install Raspberry NOOBS as detailed [here](https://projects.raspberrypi.org/en/projects/raspberry-pi-getting-started/2)
- connect the raspberry to your local WiFi
- run `sudo raspi-config` > interfacing options > SSH to enable the ssh interface
- connect via ssh `ssh pi@192.168.1.54`
- sudo apt-get install vim
- curl -fsSL get.docker.com -o get-docker.sh && sh get-docker.sh
- sudo gpasswd -a \$USER docker
- newgrp docker
- docker run hello-world
- plugin the camera
- enable camera interface sudo raspi-config > interfacing options > Camera
- raspistill -w 1600 -h 1200 --timeout 1 --brightness 50 --quality 90 --hflip null -o ~/capture.jpg

## CLI

iot-abe has a cli for interacting with the different components:

```{.sh}
make build && ./output/iot-abe

iot-abe - attribute-based access control for the IoT.

Usage:
  iot-abe [command]

Available Commands:
  client      Retrieves encrypted data from an iot device server and decrypts it (if possible)
  devices     Manages IoT devices
  grant       Manages access rights
  help        Help about any command
  request     Manages access right requests
  server      Starts an iot device server
  version     Print the version number of iot-abe
  wallet      Manages eth wallets

Flags:
  -c, --config.path string   config path (default "./config.yml")
  -h, --help                 help for iot-abe

Use "iot-abe [command] --help" for more information about a command.
```

### Examples

```{.sh}
./output/iot-abe wallet list
./output/iot-abe wallet transfer --from 0x20683Db6E6d7ff53b62BCD6F723f74eC94dC410e --to 0x1e52b030261C4890A6aCe85Ed48CaE5f459525A0 --amount 0.01
./output/iot-abe grant init --from 0xBB79396384ed533476b9D2Edf6c25797Ab3eD2cD
./output/iot-abe grant watch-requests --contract=0x7bF73B9dFA1d9A520de1Bd4BB829d4Dc602b4567
./output/iot-abe request access --contract 0x7bF73B9dFA1d9A520de1Bd4BB829d4Dc602b4567 --for 0xa9a0E7C567f5fE4f9C7f684b3398FD74041385BF
./output/iot-abe grant get-requests --contract 0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e
./output/iot-abe request watch-grants --contract 0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e
./output/iot-abe grant access --for=0x1e52b030261C4890A6aCe85Ed48CaE5f459525A0 --contract=0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e --owner=0x20683Db6E6d7ff53b62BCD6F723f74eC94dC410e --attributes="admin,ceo,it_staff"
./output/iot-abe request get-grant --for 0x1e52b030261C4890A6aCe85Ed48CaE5f459525A0 --contract 0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e
./output/iot-abe devices watch-policy-updated --contract=0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e
./output/iot-abe devices add --contract=0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e --owner=0x20683Db6E6d7ff53b62BCD6F723f74eC94dC410e --device=0xE1097bAAA914277A8E2AefE464f8E29557e5f046 --name="Camera A" --policy="(admin & it_departement)"
./output/iot-abe devices get --contract=0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e --device=0xE1097bAAA914277A8E2AefE464f8E29557e5f046
./output/iot-abe devices get-all --contract 0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e
./output/iot-abe devices watch-policy-removed --contract=0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e
./output/iot-abe devices remove --contract=0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e --device=0xE1097bAAA914277A8E2AefE464f8E29557e5f046 --owner=0x20683Db6E6d7ff53b62BCD6F723f74eC94dC410e
./output/iot-abe server
./output/iot-abe client
./output/iot-abe client serve-capture --server http://192.168.1.54:8080
```

## Help

[free space Raspberry](https://scribles.net/free-up-sd-card-space-on-raspberry-pi/)

## Links

- https://github.com/golang/go/wiki/cgo#global-functions
- https://developer.gnome.org/glib/stable/glib-Byte-Arrays.html
- https://stackoverflow.com/questions/35673161/convert-go-byte-to-a-c-char
- Have a look at: https://github.com/ethereum/EIPs/issues/1481
- https://solidity.readthedocs.io/en/develop/types.html#arrays
- https://web3js.readthedocs.io/en/1.0/web3-utils.html#hextobytes
- https://manojpramesh.github.io/solidity-cheatsheet/#dynamic-byte-arrays
- https://gist.github.com/OR13/08e2ceba147b52ef078c4527e1c48a25
- https://github.com/pubkey/eth-crypto#recoverpublickey
- https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts
- https://github.com/ethereum/go-ethereum/tree/master/crypto/ecies
- https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671
- https://ownyourbits.com/2018/06/27/running-and-building-arm-docker-containers-in-x86/
- https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671
- https://medium.com/@chrischdi/cross-compiling-go-for-raspberry-pi-dc09892dc745
- https://github.com/balena-io-library
- https://stackoverflow.com/questions/54842833/access-raspistill-pi-camera-inside-a-docker-container
- https://github.com/sgerrand/docker-glibc-builder

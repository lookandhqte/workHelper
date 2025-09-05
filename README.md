# workhelper

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

A service authomatizating default actions on hh.ru

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Install

```sh
go get github.com/lookandhqte/workHelper
```

## Usage

```sh
beanstalkd -l 127.0.0.1 -p 11300
go run cmd/app/main.go
go run cmd/worker/main.go
```

## Maintainers

[@Gelena](https://lookandhqte.com)




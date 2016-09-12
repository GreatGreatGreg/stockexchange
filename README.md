# StockExchange

Stockexchange is a web platform for stock trading. You can visit its live
instance [here](https://stockex.herokuapp.com/).

# Getting Started

In order to build and run the application you will need the recent version of Go 1.7.1.

```sh
$ git clone https://github.com/svett/stockexchange
$ git submodule update --init --recursive
$ go run main.go
```

#### Running tests

In order to start contributing to the project, you should install
[ginkgo](http://github.com/onsi/ginkgo) and
[gomega](http://github.com/ons/gomega) package that are used in unit and
integration tests:

```bash
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega
```

```bash
# Running the integration tests
$ ginkgo integration/
# Running the unit tests
$ ginkgo .
```

# License

*MIT*

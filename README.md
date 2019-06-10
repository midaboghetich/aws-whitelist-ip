# aws-withelist-ip

<h2 align="center">A Simple GO application to whitelist an IP address in AWS</h2>

## Pre-requisites
you should have a valid [AWS configuration and credential files](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)
## Getting Started

To whitelist your current ip address just run (the app will use https://api.ipify.org to retrieve your ip address):

`./aws-whitelist-ip -g <the security group to modify> -email <your email address> -p <the port number to expose>`

To whitelist a specific IP address run: 

`./aws-whitelist-ip -g <the security group to modify> -ip <the ip address to whitelist> -email <your email address> -p <the port number to expose>`

## Bugs, feature requests and code improvement
If you have a bug, a feature request or you think that the code could be improved (this is my first attempt to code in GO), please use the [issue tracker](https://github.com/midaboghetich/aws-whitelist-ip/issues).

## Creator

**Mida Boghetich**

- <http://www.boghetich.com>
- <https://github.com/midaboghetich>

## Copyright and License

Code Copyright 2019 Mida Boghetich. Code released under the [Apache License 2.0](https://github.com/midaboghetich/aws-whitelist-ip/blob/master/LICENSE).

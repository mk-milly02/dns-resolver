# DNS Resolver

## Overview

The DNS Resolver is a project designed to resolve domain names to their corresponding IP addresses.
This tool can be used to perform DNS lookups and retrieve various DNS records such as A, AAAA, CNAME, MX, and more.

## Features

- Resolve domain names to IPv4 (A) and IPv6 (AAAA) addresses.
- Retrieve A, NS, AAAA, and other DNS records.
- Support for both synchronous and asynchronous operations.
- Command-line interface for easy usage.

## Usage

To install the DNS Resolver, clone the repository and install the required dependencies:

```bash
git clone https://github.com/mk-milly02/dns-resolver.git
cd dns-resolver
go run main.go -d "example.com"
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.
Ensure that your code follows the project's coding standards and includes appropriate tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

<!-- SPDX-License-Identifier: MIT OR Apache-2.0 -->

# GO ECS
[![License](https://img.shields.io/badge/license-MIT%2FApache--2.0-informational?style=flat-square)](COPYRIGHT.md)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

Go ECS is a simple sparse set, query based Entity Component System (ECS) library.

## Description

This ECS implementation uses a sparse set style data structure for packing the
entity and component (arbitrary data) arrays while allowing O(1) entity lookups.
It is simple and lightweight.

Iterating over a list of entities and components is optimal. However, iterating
over archetypes with multiple components is only supported through O(N) queries.
There is no grouping feature to improve query time as of now. However, this is
still very fast and allows for extremely fast add and remove operations.

Creation and destruction must be handled by the user. Systems are not managed
by the world: there is no scheduler or event system. There are only queries.
The rest is up to the programmer. This is not a framework, just another tool.

This package has no external dependencies and avoids reflection by way of Go's
limited generic types. As a tradeoff, a separate query function must be written
for every number of components queried (Go does not support variadic generics).

## Usage

Go get the package.

```zsh
go get github.com/jdavasligil/go-ecs@latest
```

Import the package.

```golang
import (
    "github.com/jdavasligil/go-ecs"
)
```

See the examples folder for code examples. The basic example
should be enough to get you going.

You can run the example with the following.

```zsh
go run ./examples/basic
```

## Authors

[J. Davasligil](jdavasligil.swimming625@slmails.com)

## Version History

* 1.0.0

## Contributing
Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall be
dual licensed as below, without any additional terms or conditions.

## License

&copy; \<2024\> \<Jaedin Davasligil\>.

This project is licensed under either of

- [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE-APACHE`](LICENSE-APACHE))
- [MIT license](https://opensource.org/licenses/MIT) ([`LICENSE-MIT`](LICENSE-MIT))

The [SPDX](https://spdx.dev) license identifier for this project is `MIT OR Apache-2.0`.

## Acknowledgments

* [Austin Morlan](https://austinmorlan.com/posts/entity_component_system/#what-is-an-ecs)
* [Dakom](https://gist.github.com/dakom/82551fff5d2b843cbe1601bbaff2acbf)
* [Michele Caini](https://skypjack.github.io/2019-02-14-ecs-baf-part-1/)
* [awesome-readme](https://github.com/matiassingers/awesome-readme)
* [project-layout](https://github.com/golang-standards/project-layout)

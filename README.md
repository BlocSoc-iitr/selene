![build](https://github.com/BlocSoc-iitr/selene/actions/workflows/go.yml/badge.svg)
![tests](https://github.com/BlocSoc-iitr/selene/actions/workflows/test.yml/badge.svg)
![linter](https://github.com/BlocSoc-iitr/selene/actions/workflows/cilint.yml/badge.svg)


# experiment

NOTE: The functions need not have correct parameters and return data types yet


## Suggestions

- `Option<u8>` can be implemented by taking a pointer to a `uint8` value ?
- we can create a `constants.go` if needed.
- if we need to modify some struct datatype , we can have a method taking it as a pointer and modify , instead of returning after making a new copy

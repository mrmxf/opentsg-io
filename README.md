# opentsg-io

## Description

opentsg-io contains the  encoding methods for custom image libraries, which are
not part of the golang/image library. This module is part of open tpg project,
but can be used independently for when you just need an image saved.

Thescurrent image files are:

- DPX files.
- Tiff files, which are part of the golang library, but this saves without any
  alpha channel.
- EXR files, there are alternatives avaiable.
- CSV representing the red blue green channels - MAY remove

At the moment no decoders are contained within in the repo, but this much change
with the needs of the project.

## Limtations

These functions are made to work with opentsg, so often use *image.NRGBA64
instead of image.Image, because we do not handle the other images. If you would
like to contribute to this to improve the options available to people who would
like to use this package please do.

**DPX:**

- 8 12 and 16 bits
- methods of encoding. little endian

**EXR:**

- only one method etc

## Visuals

Depending on what you are making, it can be a good idea to include screenshots
or even a video (you'll frequently see GIFs rather than actual videos). Tools
like ttygif can help, but check out Asciinema for a more sophisticated method.

## Installation

```sh
go get opentsg-io
```

## Usage

```go
import (
    image/image
    image/draw
    github.com/mrmxf/opentsg-io/dpx
)

func main() {

// make an image

// paint it green

//save it as a dpx


}
```

## Support

Raise an issue at [https://github.com/mrmxf/opentsg-io].

## Roadmap

If you have ideas for releases in the future, it is a good idea to list them in the README.

## Contributing

State if you are open to contributions and what your requirements are for accepting them.

For people who want to make changes to your project, it's helpful to have some documentation on how to get started. Perhaps there is a script that they should run or some environment variables that they need to set. Make these steps explicit. These instructions could also be useful to your future self.

You can also document commands to lint the code or run tests. These steps help to ensure high code quality and reduce the likelihood that the changes inadvertently break something. Having instructions for running tests is especially helpful if it requires external setup, such as starting a Selenium server for testing in a browser.

## Authors and acknowledgment

Show your appreciation to those who have contributed to the project.

## License

BSD-3 clause with attribution

## Project status

If you have run out of energy or time for your project, put a note at the top of the README saying that development has slowed down or stopped completely. Someone may choose to fork your project or volunteer to step in as a maintainer or owner, allowing your project to keep going. You can also make an explicit request for maintainers.

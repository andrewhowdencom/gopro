# GoPro (Linux)

This is a set of tools for interacting with your GoPro through a friendly CLI interface.

Only available on Linux, as Gopro generally doesn't support it.

## Support

This application has no community support. It is just a fun project from the author.

However, the operating systems in which this is expected to function are:

- Debian (Buster)

Additionally, it has been tested with the GoPro versions:

- HERO 8

## Compilation

The application compilation is managed by `make`, through a self documenting makefile.

Run `make` from the project root to see the help menu.

### Just give me the .deb

To compile the deb, run:

```bash
make clean app debn
```

Then, to install it run:

```
sudo deb.install
```


# Kody

Command line helper for Epic React Dev.

The CLI is named after Kody, the mascot for the Epic React course!

## Build

Currently, the only way to run the project is to have Go 1.23+ installed locally and building a project.

```bash
$ go build
$ ./kody --help # or ./kody.exe on Windows
```

## Commands

### Save

Save takes the current exercise in the Playground and saves it to a more permanent location:

```bash
$ kody save --workshop react-fundamentals --output ~/epic-react-solutions
```

This command goes into the `react-fundamentals` directory, and copies the contents of the playground to a folder inside `~/epic-react-solutions`.
The use-case for this command is to save the playground's contents of an exercise you just finished into a more permanent location, like a git repository.

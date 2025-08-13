# Kody

Command line helper for [Epic React Dev](https://www.epicreact.dev/) course by Kent C. Dodds.

Kody saves the exercise solutions from your playground if you later want to check them out. Useful if you are reviewing the exercises you've made.

Kody is also able to restore previous solutions to the currrent playground. This is useful if you a revisitng an older exercise and want to bring back the solution you made to the playground.

The CLI is named after Kody, the mascot for the Epic React course!

## Installation

### Option 1: Download from releases

**Download Links for v1.0.0:**

- [Windows](https://github.com/andrerfcsantos/kody/releases/download/v1.0.0/kody_Windows_x86_64.zip) 
- Mac ([Intel](https://github.com/andrerfcsantos/kody/releases/download/v1.0.0/kody_Darwin_x86_64.zip) | [Apple](https://github.com/andrerfcsantos/kody/releases/download/v1.0.0/kody_Darwin_arm64.zip))
- Linux ([x86_64](https://github.com/andrerfcsantos/kody/releases/download/v1.0.0/kody_Linux_x86_64.tar.gz) | [arm64](https://github.com/andrerfcsantos/kody/releases/download/v1.0.0/kody_Linux_arm64.tar.gz))

See the full list of versions of the binaries for diferent OS/architctures in the [GitHub releases page](https://github.com/andrerfcsantos/kody/releases).

All releases come with a binary that you can run directly.

Add it to your PATH for easier access to the `kody` command anywhere in the system!

### Option 2: Build from source

For the following options, you should have Go 1.23+ installed locally.
Read more on compiling and installing a Go application [in this tutorial](https://go.dev/doc/tutorial/compile-install).

#### Install directly from remote repository
```bash
go install github.com/andrerfcsantos/kody@latest
```

#### OR: Clone and build locally
```bash
git clone https://github.com/andrerfcsantos/kody.git
cd kody
go build # Or 'go install' to make the binary available globally 
./kody --help # or ./kody.exe on Windows
```

## Typical workflow example with kody

This is a typical use of kody, where you save your progress as you complete exercises and then restore previous solutions when needed. 

```bash
# 1. Configure Kody for your setup. You only need to do this once.
kody config workshops.dir ~/epic-react-workshops
kody config save.output.directory ~/epic-react-solutions

# 2. Check current exercise status
kody status

# 3. You solve the exercise in the playground and run the tests on the workshop until they pass

# 4. Save your completed exercise once you are done
kody save

# 5. You do more exercises, but then come back to a previous one and click "set playground to this exercise" in the workshop

# 6. Restore the solution you had for this exercise
kody restore
```

## Configuration

Before using Kody, it's recommended you set up some initial configuration.
That way you don't have to repeat yourself when giving commands to kody!

The two recommended settings are:

- `workshops.dir`: Path to the directory containing all workshop subdirectories. The current workshop and exercise you are working on will be auto-detected if all your workshops are in this folder.
- `save.output.directory`: (Optional) Directory where exercises will be saved. This can be a git repository where you are saving all your solutions for the exercises. If you don't pass a value to this configuration, Kody will save the exercises to a data folder on your machine. No need to configure this if you don't care too much where the exercises are being stored.

### Recommended configuration

```bash
# Set the workshops directory (should contain all Epic React workshops in subdirectories)
kody config workshops.dir ~/epic-react-workshops

# Set the output directory for saving exercise solutions (or skip this to have them saved in a default location)
kody config save.output.directory ~/epic-react-solutions
```

Replace `~/epic-react-workshops` and `~/epic-react-solutions` with directories that make sense for you.

### Advanced configuration

#### Auto-commit

If your `save.output.directory` is also a git repository, you can ask kody to auto-commit for you everytime you save an exercise:

```
kody config save.shouldCommit true
```

#### Opt-out of workshop auto-detection

If you don't want the workshop to be auto-detected with `workshops.dir`, you can specify a workshop folder with the workshop you are currently working:

```
kody config workshop.dir ~/epic-react-workshops/react-fundamentals
```

You then must run this command again every time you change workshops.


## Commands

### Save

Save takes the current exercise in the Playground in the workshop you are currently working on, and saves it to a more permanent location.

#### Simple usage (assumes previous configuration)

Provided you've configurated the workshops folder and a directory where to save exercises to you can simply run this to save an exercise:

```bash
# Save current exercise from the current workshop (both auto-detected!)
kody save
```

Kody will try to figure out the workshop you are working on based on the one that is inside `workshops.dir` with the most recent modification time.
Inside that workshop, it'll also auto-detect the exercise you are doing based on the contents of the playground.
Then, it'll take the playground folder contents, and save it to the `save.output.directory` you've previously configured.

#### Custom usage with flags

You can also pass flags to override the configuration you've previously set up or to specify things you didn't setup a config for:

```bash
# Save current exercise of a specific workshop to a directory
kody save --workshop ~/epic-react-workshops/react-fundamentals --output ~/my-solutions

# Save with auto-detection of current workshop
kody save --workshops-dir ~/epic-react-workshops

# Save and commit changes to git
kody save --commit

# Use short flags
kody save -w ~/epic-react-workshops/react-fundamentals -o ~/my-solutions -c
```

### Restore

Restore an exercise to the playground from a previously saved location.

#### Simple usage (assumes previous configuration)
```bash
# Restore current exercise automatically
kody restore
```

Kody will auto-detect the exercise and workshop you are working on, and will fetch your previously saved solution for that exercise and place it in the workshop folder for you.

#### Custom usage with flags
```bash
# Restore specific exercise by section and exercise number
kody restore 01.02

# Restore with custom workshop path
kody restore 01.02 --workshop ~/epic-react-workshops/react-fundamentals

# Restore with auto-detection of current workshop
kody restore 01.02 --workshops-dir ~/epic-react-workshops

# Use short flags
kody restore 01.02 -w ~/epic-react-workshops/react-fundamentals
```

### Status

Get information about the current workshop and exercise and based on the contents of the playground.
Use this command to double-check kody is correctly detecting the workshop/exercise you are working on.

#### Simple usage (assumes previous configuration)
```bash
# Show current exercise status
kody status
```

#### Custom usage with flags
```bash
# Check status of specific workshop
kody status --workshop ~/epic-react-workshops/react-fundamentals

# Check status with auto-detection
kody status --workshops-dir ~/epic-react-workshops

# Use short flags
kody status -w ~/epic-react-workshops/react-fundamentals
```

### Config

Manage Kody configuration settings.

```bash
# View all configurations
kody config

# View specific configuration value
kody config workshops.dir

# Set configuration value
kody config workshops.dir ~/epic-react-workshops
```

### Version

Display version information.

```bash
kody version
```


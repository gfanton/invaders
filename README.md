# 👽 Alien Invaders
Alien Invaders is a simulation game implemented in Go. 
The game involves aliens who randomly traverse a map filled with cities. 
Whenever two aliens encounter each other in the same city, they engage in a battle, leading to their mutual destruction and the annihilation of the city they were in. 
The simulation ends when all aliens have been killed or when aliens become exhausted.

## 🛠️ Installation
1. Ensure that Go (v1.19 or later) is installed on your system.
2. Run `make` in the project root directory to install the game.

**Note**: You can modify the `INSTALL_DIR` variable in the Makefile to change the installation directory. The default directory is `/usr/local/bin`.

## 🎯 Usage

```bash
invader [flags] <subcommand>
```

### Subcommands

#### 1. `start`
This subcommand initiates the Alien Invaders simulation. The
program reads from standard input by default, but you can
specify a file instead.

```bash
USAGE
  invader start -alien [value] -file [path] -max_steps [value]

FLAGS
  -alien 10         The seed used to start the map; an empty seed will be selected if not provided.
  -file string      Read from a specified file instead of the standard input.
  -max_steps 10000  The maximum number of steps an alien can perform before becoming exhausted.
```

#### 2. `generate`
This subcommand is used to generate a new random city map of a given depth.

```bash
USAGE
  invader generate -depth [value] -seed [string]

FLAGS
  -depth 100    the depth of the wanted map
  -seed string  the seed used to generate the map, empty seed will be choose if empty
```

## Example
a fast way to test this program is to cumulate generate + start:

```bash
 invader generate -depth=5 | tee /dev/tty | invader start -aliens=4 
```

## Running Tests
To run tests, navigate to the project directory and run the following command:

```bash
make test
```

## TODO
These cool enhancements could be made when time permits:

* [ ] Implement a graph representation: Visualizing the movements of the aliens
      across the map could be fascinating, but it might be challenging to implement
      or visualize in the terminal. A JavaScript interface may be more suitable for
      this.
* [ ] Use different level of verbosity: for now, debug lecture can be difficult
      with bigger map, multiple verbosity level should help imporove debuging
* [ ] Decouple printing logic: It might be beneficial to decouple printing and
      logging functionality from the simulation by using separate structures or
      interfaces.
  
## Notes
* An alternative approach to this exercise could involve making the aliens
  autonomous using goroutines and syncing them with some sort of ticker. This
  could potentially simplify the management of their lifecycle, but add syncing
  complexity that's not really necessary
* I voluntary add "trapped" message that should be mendatory for me, every *non mendatory* 
  messages should appears in the debug mode


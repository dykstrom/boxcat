# BoxCat

Think inside the box. With cats.

BoxCat is a whimsical, cat-themed esoteric programming language where all computation happens with the help of a curious Cat and a collection of cardboard boxes.

![A cat in a box](boxcat.png "Think inside the box. With cats.")

## What is BoxCat?

BoxCat is a programming language inspired by the unpredictable and playful nature of cats. Each program consists of commands that direct a Cat to interact with boxes, move values around, perform calculations, and control the flow of execution.

- **Boxes** are named memory cells that can hold integer values.
- **The Cat** acts as the accumulator, holding a single integer value in its Paw.
- **Commands** are written in a readable, cat-themed syntax.

## Example: Print the numbers 1 to 5

```text
POUNCE ON 1         # Put 1 in the Cat's Paw
SIT IN COUNTBOX     # Store 1 in COUNTBOX, Paw is now 0

LOOP_START:
JUMP OUT OF COUNTBOX    # Paw = COUNTBOX, COUNTBOX = 0
MEOW                    # Output the current number
POUNCE ON 1
PURR AT COUNTBOX        # Paw = Paw + COUNTBOX (increment)
SIT IN COUNTBOX         # Store incremented value back in COUNTBOX, Paw = 0
PEEK INSIDE COUNTBOX    # Paw = COUNTBOX
POUNCE ON 5
HISS AT COUNTBOX        # Paw = 5 - COUNTBOX
IF CAT CURIOUS, LEAP TO LOOP_START  # If not zero, loop again

TAKE A NAP              # End program
```

This program prints the numbers 1 through 5, each on its own line.

_Note: Does this program work? COUNTBOX will always be 0 at the start of the loop, because the cat jumps out of COUNTBOX. Better to peek inside COUNTBOX maybe?_

## Learn More

See the [language specification](specification.md) for the full list of commands.

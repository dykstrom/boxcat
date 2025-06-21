# BoxCat Language Specification

## Core Concept

The program state is represented by a collection of cardboard boxes and a single Cat. The Cat can move between boxes, manipulate items (numeric values) within boxes, and perform actions based on the contents or state of the boxes.

## Memory Model

- **Boxes:** A collection of named cardboard boxes. Each box can hold a single integer value. Boxes are created implicitly when first referenced. If a box is empty (has not been `SIT IN` or was the source of a `JUMP OUT OF` command without being refilled), its value is considered 0.
- **The Cat (Accumulator):** There is one Cat. The Cat always holds a single integer value, referred to as what's "in its Paw." The Cat starts with 0 in its Paw.

## Syntax

- Programs are a sequence of commands, one per line.
- Commands are case-insensitive (e.g., `SIT IN` is the same as `sit in`).
- Anything after a `#` (hash symbol) on a line is a comment and is ignored.
- Box names can be any alphanumeric string (e.g., `BIGBOX`, `SECRETBOX`, `BOX1`).
- Labels for jumps are defined by ending a line with a colon, e.g., `LOOP_START:`.

## Commands

### Movement & Box Interaction

- `SIT IN [BOX_NAME]`
  - The Cat moves to the specified box.
  - The value currently in the Cat's Paw is placed into `[BOX_NAME]`, overwriting its previous contents.
  - The Cat's Paw becomes empty (its value is set to 0).
- `JUMP OUT OF [BOX_NAME]`
  - The Cat takes the value from `[BOX_NAME]` into its Paw.
  - `[BOX_NAME]` becomes empty (its value is set to 0).
  - The Cat is now "outside" any particular box, but still holds the value in its Paw.
- `PEEK INSIDE [BOX_NAME]`
  - The Cat copies the value from `[BOX_NAME]` into its Paw, overwriting what was previously in its Paw.
  - The contents of `[BOX_NAME]` are unchanged.
  - The Cat doesn't "move" into the box in terms of where it's considered to be "sitting."
- `DROP IN [BOX_NAME]`
  - The Cat drops a copy of the value in its Paw into `[BOX_NAME]`, overwriting its previous contents.
  - The value in the Paw is unchanged.
  - The Cat doesn't "move" into the box in terms of where it's considered to be "sitting."

### Cat Actions (Data Manipulation - Acts on Paw and Box)

- `POUNCE ON [VALUE]`
  - `[VALUE]` may be an integer or a single ASCII character enclosed in `'` characters.
  - Loads the literal integer `[VALUE]` or the ASCII code of the character directly into the Cat's Paw, overwriting its current contents.
- `PURR AT [BOX_NAME]`
  - Adds the value in `[BOX_NAME]` to the value in the Cat's Paw. The result is stored in the Cat's Paw.
  - `[BOX_NAME]`'s content is unchanged.
- `HISS AT [BOX_NAME]`
  - Subtracts the value in `[BOX_NAME]` from the value in the Cat's Paw. The result is stored in the Cat's Paw.
  - `[BOX_NAME]`'s content is unchanged.
- `PLAYFULLY BAT [BOX_NAME]`
  - Multiplies the value in the Cat's Paw by the value in `[BOX_NAME]`. The result is stored in the Cat's Paw.
  - `[BOX_NAME]`'s content is unchanged.
- `KNOCK OVER [BOX_NAME]` (Integer Division)
  - Divides the value in the Cat's Paw by the value in `[BOX_NAME]`. The integer result (truncating any remainder) is stored in the Cat's Paw.
  - `[BOX_NAME]`'s content is unchanged.
  - Division by zero makes the Cat very confused (program halts with an "Angry Cat" error).
- `LEAVE A "GIFT" IN [BOX_NAME]` (Modulo)
  - Calculates `(Value in Paw) % (Value in [BOX_NAME])`. The result is stored in the Cat's Paw.
  - `[BOX_NAME]`'s content is unchanged.
  - Modulo with zero makes the Cat very confused (program halts with an "Angry Cat" error).

### Input/Output

- `MEOW`
  - Outputs the integer value currently in the Cat's Paw to standard output, followed by a blank character.
- `YOWL`
  - Outputs the ASCII character corresponding to the integer value in the Cat's Paw to standard output (no automatic newline). If the value is outside the printable ASCII range, behavior is implementation-defined (e.g., output nothing, output a '?', or output the character if the system supports it).
- `LISTEN FOR WHISPER`
  - Reads an integer from standard input and stores it in the Cat's Paw, overwriting its current contents.
- `SNIFF AROUND`
  - Reads a single character from standard input, converts it to its ASCII integer value, and stores it in the Cat's Paw, overwriting its current contents.

### Control Flow

- `[LABEL_NAME]:`
  - Defines a label at this point in the code. Does not perform any action itself.
- `DART TO [LABEL_NAME]`
  - Subroutine call: Saves the current execution position and jumps to the command immediately following `[LABEL_NAME]`. Execution can return to the saved position using `DART BACK`.
- `DART BACK`
  - Returns from a subroutine call made by `DART TO`, resuming execution at the command following the original `DART TO`.
- `LEAP TO [LABEL_NAME]`
  - Unconditional jump: Jumps program execution to the command immediately following `[LABEL_NAME]`. Does not save the current position (no return).
- Conditional Execution:

  - `IF CAT CURIOUS, [COMMAND]`
    - If the value in the Cat's Paw is NOT zero, executes the single-line `[COMMAND]`. Otherwise, execution continues to the next command.
  - `IF CAT BORED, [COMMAND]`
    - If the value in the Cat's Paw IS zero, executes the single-line `[COMMAND]`. Otherwise, execution continues to the next command.
  - `IF BOX EMPTY [BOX_NAME], [COMMAND]`
    - If `[BOX_NAME]` is empty (contains 0), executes the single-line `[COMMAND]`. Otherwise, execution continues to the next command.
  - `IF BOX NOT EMPTY [BOX_NAME], [COMMAND]`
    - If `[BOX_NAME]` is not empty (not 0), executes the single-line `[COMMAND]`. Otherwise, execution continues to the next command.

  _Note: The `[COMMAND]` in conditional execution can be any valid single-line BoxCat command, including control flow commands such as `DART TO`, `LEAP TO`, or any other command._

### Quirky Cat Behaviors

- `SNIFF CATNIP`
  - The Cat gets a whiff of catnip!
  - A pseudo-random non-negative integer is generated and placed into the Cat's Paw, overwriting its current contents. The range of this random number would be implementation-defined (e.g., 0 to 32767, or 0 to a system's `RAND_MAX`).
- `SUDDENLY SCRATCH`
  - The Cat is suddenly overcome with an urge to scratch!
  - The value in the Cat's Paw is set to 0 (the Cat drops whatever it was holding, effectively clearing the accumulator).
- `DOZE`
  - The Cat decides to take a short nap.
  - Program execution pauses for a duration.
  - If the value in the Cat's Paw is greater than 0, the pause duration is approximately (Value in Paw) milliseconds (or some other implementation-defined time unit, e.g., seconds if Paw value is small).
  - If the value in the Cat's Paw is 0 or less, the doze is extremely short, effectively a no-op or a minimal system delay.

### Program Termination

- `TAKE A NAP`
  - Halts the program. The Cat is satisfied and goes for a long sleep.
- `GET STUCK`
  - The Cat gets stuck in a box or a loop of thought!
  - This is conceptually an infinite loop. A typical implementation would be: `GET_STUCK_INTERNAL_LABEL: LEAP TO GET_STUCK_INTERNAL_LABEL`. The command itself implies this loop.

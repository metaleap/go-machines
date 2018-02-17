# Functional CoreLang '91

Exploratory slow-paced walk through the 1991 Peyton-Jones / Lester book [Implementing Functional Languages](http://www.cs.otago.ac.nz/cosc459/books/pjlester.pdf) — in Go, instead of the deceased '80s Haskell precursor 'Miranda'.

A minimal (aka. lacking higher-level syntactic sugars) Functional Language named *Core* is implemented in various interpreter VMs.

Lexing + parsing is from scratch, not 'by the book'. But the book was followed closely to fully implement these graph-reduction machines:

- **Template Instantiation Machine:** simpler compilation, sub-optimal execution
    > _has a bug or 2 left in it for more intricate definitions/expressions: moved on to the cooler stuff before they could be addressed_
- **G-Machine:** trickier compilation, better execution

Machines still to be implemented:

- Three-Instruction Machine
- Spineless Tagless G-Machine

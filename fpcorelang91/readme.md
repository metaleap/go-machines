# Functional CoreLang '91

Exploratory slow-paced walk through the 1991 Peyton-Jones / Lester book [Implementing Functional Languages](http://www.cs.otago.ac.nz/cosc459/books/pjlester.pdf) — in Go, instead of the deceased '80s Haskell precursor 'Miranda'.

A minimal (aka. lacking higher-level syntactic sugars) Functional Language named *Core* is implemented in various interpreter VMs.

Lexing + parsing is from scratch, not 'by the book'. But the book was followed closely to fully implement these graph-reduction machines:

- **Template Instantiation Machine:** evaluation by graph-building and traversal/reduction, no compilation phase
    > _incomplete, lacking arithmetic, conditionals, and constructors, but the prelude defs work: moved on to the cooler stuff before completion_
- **G-Machine:** somewhat involved compilation to (essentially) byte-code, better execution

Machines still to be implemented:

- Three-Instruction Machine
- Spineless Tagless G-Machine

# Functional CoreLang '91

Exploratory slow-paced walk through the 1991 Peyton-Jones / Lester book [Implementing Functional Languages](http://www.cs.otago.ac.nz/cosc459/books/pjlester.pdf) — in Go, instead of the deceased '80s Haskell precursor 'Miranda'.

A minimal (aka. lacking higher-level syntactic sugars) Functional Language named *Core* is implemented in various interpreter VMs.

> To play, you run `cmd/corelang-dummyrepl` that picks up the most essential definitions from `prelude.go` (SKI++ basically) and parses various playful extra definitions from its `dummy-mod-src.go`, also readlines ad-hoc user definitions and invokes evaluation runs when prompted.

Lexing + parsing (in `syn`) is from scratch, not 'by the book'. But the book was followed closely to implement these graph-reduction machines:

- **Template Instantiation Machine:** evaluation by graph-building and traversal/reduction, no compilation phase
    > _incomplete, lacking arithmetic, conditionals, and constructors, but the prelude-defs work: moved on to the cooler stuff before completion, given this machine's real-world uselessness (except for introductory teaching)_
- **G-Machine:** completed all levels (mark 7) — somewhat involved compilation to (essentially) byte-code, better evaluation characteristics (flat pre-generated instruction stream instead of ad-hoc graph construction/traversal)

Machines still to be implemented:

- Three-Instruction Machine
- Spineless Tagless G-Machine

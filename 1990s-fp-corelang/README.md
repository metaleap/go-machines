# Early-90s functional language machines

Exploratory slow-paced walk through:
- the 1991 Peyton-Jones / Lester book [Implementing Functional Languages](http://www.cs.otago.ac.nz/cosc459/books/pjlester.pdf) — in Go, instead of the deceased '80s Haskell precursor 'Miranda',
- the 1992 Peyton-Jones [Spineless Tagless G-Machine](https://www.microsoft.com/en-us/research/wp-content/uploads/1992/04/spineless-tagless-gmachine.pdf) utilizing/re-purposing this repo's now (as per above) existing (parsing/compiling/VM etc.) machinery

A minimal (aka. lacking higher-level syntactic sugars) Functional Language named *Core* is implemented in various interpreter VMs.

> To play, you run `cmd/corelang-dummyrepl` that picks up the most essential definitions from `prelude.go` (SKI++ basically) and parses various playful extra definitions from its `dummy-mod-src.go`, also readlines ad-hoc user definitions and invokes evaluation runs when prompted.

Lexing + parsing (in `syn`) is from scratch, not 'by the book'. But the above materials were then followed closely to approach implementing these graph-reduction machines as follows:

- **Template Instantiation Machine:** evaluation by ad-hoc graph construction/instantiation and traversal/reduction, no real separate pre-processing / compilation stage
    > _incomplete: lacking (functioning) arithmetic, conditionals, and constructors, but the prelude-defs work — had to move on to the cooler stuff below, given the template-instantiation machine's real-world uselessness (except for newcomers testing the waters)_
- **G-Machine:** completed all levels (mark 7 — but somehow still failed to do proper lambdas / lambda-lifting / non-nilary local LET defs — still, one must move on) — involves fairly intricate compilation schemes to this virtual reduction-machine's rather convoluted (essentially) byte-code — better (than above one, but still atrocious) execution characteristics (linear pre-generated instruction stream instead of graph traversal)
- **Three-Instruction Machine:** — skipped entirely
- **Spineless Tagless G-Machine:** — in progress

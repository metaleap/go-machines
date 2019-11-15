main args   := if ?args (factorial (num.parseUnsigned (args.at 0))) (ERR errmsg)
  errmsg    := "missing: int arg to factorial"

factorial i := if (i.EQ 0) 1 (else i.MUL (factorial (i.SUB 1)) )

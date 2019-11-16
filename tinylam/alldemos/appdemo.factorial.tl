main args   := if (! args) (ERR errmsg) (num.factorial (num.parseUnsigned (args.-<)))
  errmsg    := "missing: int arg to factorial"

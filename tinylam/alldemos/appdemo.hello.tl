main args := ** concat ["Hello ", name, "!"]
  name    := if (! args) "World" (else args.-<)

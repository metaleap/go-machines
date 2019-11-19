main args := concat.** ["Hello ", name, "!"]
  name    := args ? [] => "World" | .. => ever

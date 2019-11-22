main args env := concat.** ["Hello ", name, "!"]
  name        := args ? [] => "World" | .. => ever

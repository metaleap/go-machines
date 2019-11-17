main args := concat.** ["Hello ", name, "!"]
  name    := args ? [] => "World" | => True

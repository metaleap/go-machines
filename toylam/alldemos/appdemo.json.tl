main args env stdin := jsondom Result? Yay => json.stringified | Nay => msg -> ERR msg stdin
  jsondom           := json.parsed stdin

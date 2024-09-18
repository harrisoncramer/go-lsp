# go-lsp

This is a proof-of-concept for how to develop an [LSP](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/) in Go for Neovim.

The basic functionality is taken from TJ's [awesome Youtube video](https://www.youtube.com/watch?v=YsdlcQoHqPY&ab_channel=TJDeVries) on working with LSPs, and provides some additional utilities for hot-reloading and improves the Go code.

[![Demo](https://hjc-public.s3.amazonaws.com/lsp-preview.png)](https://hjc-public.s3.amazonaws.com/lsp-demo.mp4)

# Dependencies

- [Go v1.23](https://go.dev/)
- [Air](https://github.com/air-verse/air)
- [Task](https://github.com/go-task/task)
- Neovim


# Development

1. Clone this repository to Neovim's `lua` configuration folder:

```bash
git clone git@github.com:harrisoncramer/go-lsp.git ~/.path-to-your-config/lua
```

2. Add a command to reload the server to your Neovim configuration. Provide it your LSP attach function as a callback:

```lua
vim.api.nvim_create_user_command("<leader>R", function()
  local on_attach = require("lsp.init").on_attach # Your LSP handler w/ keybindings, etc.
  require("go-lsp").restart(on_attach)
end, { nargs = 0 })
```

3. Start the Air binary:

```bash
cd ~/.path-to-your-config/lua/go-lsp
task dev
```

4. Open up a markdown file and press `<leader>R` to start the development server. When you make changes to the LSP source code, the `air` binary will rebuild the LSP automatically. Run `<leader>R` to reattach to the rebuilt LSP server.

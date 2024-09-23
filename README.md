# go-lsp

The LSP protocol defines how a server and client must exchange information when providing language features like auto-completion, go-to-definition, or diagnostics. This repository aims to simplify the process of creating an LSP server in Go, tailored specifically for integration with Neovim.

It contains a basic implementation of the LSP protocol, written in Go. It also provides a Lua script and tooling for hot-reloading the server in Neovim.

[![Demo](https://hjc-public.s3.amazonaws.com/lsp-preview.png)](https://hjc-public.s3.amazonaws.com/lsp-demo.mp4)

It's meant to server as a jumping off point for quickly developing real LSP functionality. The basic code structure is cribbed from TJ's [awesome Youtube video](https://www.youtube.com/watch?v=YsdlcQoHqPY&ab_channel=TJDeVries) on working with LSPs.

## Dependencies

- [Go v1.23](https://go.dev/)
- [Air](https://github.com/air-verse/air)
- [Task](https://github.com/go-task/task)
- Neovim

## Development

1. Clone this repository to Neovim's `lua` configuration folder:

```bash
git clone git@github.com:harrisoncramer/go-lsp.git ~/.path-to-your-config/lua/
```

2. Add a command to reload the server to your Neovim configuration. Provide it your LSP attach function as a callback, so you can set up keybindings and so forth:

```lua
-- Or whatever your "on_attach" function looks like...
local on_attach = function(client, bufnr)
  vim.keymap.set("n", "K", vim.lsp.buf.hover, { noremap = true, silent = true })
end

vim.keymap.set("n", "<leader>R", function()
  local ok, go_lsp = pcall(require, "go-lsp")
  if not ok then
    error("Could not load Go LSP!")
  end
  go_lsp.restart(on_attach)
end)
```

3. Start the Air binary:

```bash
cd ~/.path-to-your-config/lua/go-lsp
task dev
```

4. Open up a markdown file and press `<leader>R` to start the development server. When you make changes to the LSP source code, the `air` binary will rebuild the LSP automatically. Run `<leader>R` to reattach to the rebuilt LSP server.

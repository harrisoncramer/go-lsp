local M = {
  existing_client = 0,
  cmd_id = 0
}

---Reattaches to the new Language Server
---@param on_attach function
M.restart = function(on_attach)
  -- If a previous autocmd exists, delete it
  if M.cmd_id ~= 0 then
    vim.api.nvim_del_autocmd(M.cmd_id)
  end

  -- If a previous client exists, stop it
  if M.existing_client ~= 0 then
    vim.lsp.buf_detach_client(0, M.existing_client)
    vim.lsp.stop_client(M.existing_client)
  end

  -- Get this directory...
  local current_file = debug.getinfo(1, "S").source:sub(2)
  local current_dir = current_file:match("(.*/)")
  local bin_path = current_dir .. "/tmp/main"

  if not vim.loop.fs_stat(bin_path) then
    error("Must build the binary first!")
  end

  -- Start the new server
  vim.lsp.start_client({
    name = "go-lsp",
    cmd = { bin_path },
    on_attach = on_attach,
    on_init = function(client)
      vim.notify("Client ready...")
      M.existing_client = client.id
      M.cmd_id = vim.api.nvim_create_autocmd("FileType", {
        pattern = "markdown",
        callback = function()
          vim.lsp.buf_attach_client(0, client.id)
        end,
      })

      local current_buf = vim.api.nvim_get_current_buf()
      local current_filetype = vim.bo[current_buf].filetype
      if current_filetype == "markdown" then
        vim.lsp.buf_attach_client(current_buf, client.id)
      end
    end,
  })
end

return M

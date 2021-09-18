set spell
set spelllang=en
set textwidth=79
set colorcolumn=80
let g:goyo_width = 103

autocmd FileType markdown setlocal ts=2 sts=2 sw=2 expandtab textwidth=99 colorcolumn=100
autocmd FileType markdown setlocal nofoldenable
autocmd BufRead,BufNewFile *.md setlocal spell wrap

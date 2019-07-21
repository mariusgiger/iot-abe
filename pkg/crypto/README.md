# Bugs

- Makefile of linux does not add libgmp at the right place in the linker command of libcpabe use `Makefile.linux` to fix this behaviour.
- policy_lang.y is missing a semicolon. fix it by changing line 67 to `result: policy { final_policy = $1; }`

## Links

- https://ubuntuforums.org/showthread.php?t=2254939

# Overriding Test Configuration

To override any test configuration, we merge multiple files into a single struct.

You can specify multiple file paths using `CTF_CONFIGS=path1,path2,path3`.

The framework will apply these configurations from right to left.

> [!NOTE]  
> When override slices remember that you should replace the full slice, it won't be extended by default!


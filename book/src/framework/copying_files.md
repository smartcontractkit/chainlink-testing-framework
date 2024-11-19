# Copying Files

You can copy files to containers by using the `ContainerName` from the output and specifying the source `src` and destination `dst` paths.

However, using this API is discouraged and will be **deprecated** in the future, as it violates the principles of "black-box" testing. If your service relies on this functionality, consider designing a configuration or API to address the requirement instead.

```go
	bc, err := blockchain.NewBlockchainNetwork(&blockchain.Input{
	    ...
	})
	require.NoError(t, err)

	err = dockerClient.CopyFile(bc.ContainerName, "local_file.txt", "/home")
	require.NoError(t, err)
```

{ stdenv, pkgs, lib }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    foundry-bin
    go-ethereum
    solc-select
    go
    gopls
    delve
    golangci-lint
    gotools
    jq
  ];
  GOROOT="${pkgs.go}/share/go";

  shellHook = ''
  solc-select install 0.8.19
  solc-select use 0.8.19
  '';
}

{ stdenv, pkgs, lib }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    foundry-bin
    go-ethereum
    solc-select
    go_1_21
    gopls
    delve
    gotools
    jq
  ];
  GOROOT="${pkgs.go}/share/go";

  shellHook = ''
  solc-select install 0.8.19
  solc-select use 0.8.19
  # setup go bin for nix
  export GOBIN=$HOME/.nix-go/bin
  mkdir -p $GOBIN
  export PATH=$GOBIN:$PATH
  # workaround to install newer golangci-lint, so that we don't have to update all dependencies
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
  '';
}

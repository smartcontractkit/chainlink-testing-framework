{ stdenv, pkgs, lib }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gopls
    delve
    golangci-lint
    gotools
    jq
  ];
  GOROOT="${pkgs.go}/share/go";

  shellHook = ''
    export PATH=$GOPATH/bin:$PATH
    go install cmd/havoc.go
  '';
}

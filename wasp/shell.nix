{ stdenv, pkgs, lib }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gopls
    delve
    golangci-lint
    gotools
    kubectl
    kubernetes-helm
    jq
  ];
  GOROOT="${pkgs.go}/share/go";

  shellHook = ''
  '';
}

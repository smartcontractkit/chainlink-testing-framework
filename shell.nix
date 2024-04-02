{ pkgs }:
with pkgs;
let
  go = pkgs.go_1_21;
  postgresql = postgresql_14;
  nodejs = nodejs-18_x;
  nodePackages = pkgs.nodePackages.override { inherit nodejs; };
in
mkShell {
  nativeBuildInputs = [
    go
    goreleaser
    postgresql

    python3
    python3Packages.pip

    curl
    nodejs
    nodePackages.pnpm
    pre-commit
    go-ethereum # geth
    go-mockery

    # tooling
    gotools
    gopls
    delve
    golangci-lint
    github-cli
    jq

    # deployment
    awscli2
    devspace
    kubectl
    kubernetes-helm
    k9s
  ] ++ lib.optionals stdenv.isLinux [
    # some dependencies needed for node-gyp on pnpm install
    pkg-config
    libudev-zero
    libusb1
  ];
  LD_LIBRARY_PATH = "${stdenv.cc.cc.lib}/lib64:$LD_LIBRARY_PATH";
  GOROOT = "${go}/share/go";

  shellHook = ''
    pre-commit install
    # Setup helm repositories
    helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo update
  '';
}

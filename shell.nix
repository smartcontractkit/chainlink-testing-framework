{ pkgs }:
with pkgs;
let
  go = pkgs.go_1_21;
  postgresql = postgresql_15;
  nodejs = nodejs-18_x;
  nodePackages = pkgs.nodePackages.override { inherit nodejs; };
  isDarwin = pkgs.stdenv.isDarwin;
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
    nodePackages.yarn
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

  LD_LIBRARY_PATH = lib.makeLibraryPath [pkgs.zlib stdenv.cc.cc.lib]; # lib64
  GOROOT = "${go}/share/go";

  shellHook = ''
    # disable CGO by default
    export CGO_ENABLED=0
    # enable pre-commit hooks
    pre-commit install
    # Setup helm repositories
    helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo update
  '';
}

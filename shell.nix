{ pkgs, scriptDir }:
with pkgs;
let
  go = pkgs.go_1_21;
  postgresql = postgresql_15;
  nodejs = nodejs-18_x;
  nodePackages = pkgs.nodePackages.override { inherit nodejs; };

  mkShell' = mkShell.override {
    # The current nix default sdk for macOS fails to compile go projects, so we use a newer one for now.
    stdenv = if stdenv.isDarwin then overrideSDK stdenv "11.0" else stdenv;
  };
in
mkShell' {
  nativeBuildInputs = [
    git
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
    dasel
    typos

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
  CGO_ENABLED = "0";
  HELM_REPOSITORY_CONFIG = "${scriptDir}/.helm-repositories.yaml";

  shellHook = ''
    # enable pre-commit hooks
    pre-commit install > /dev/null
    # enable pre-push hooks
    pre-commit install -f --hook-type pre-push > /dev/null
    # Update helm repositories
    helm repo update > /dev/null
    # setup go bin for nix
    export GOBIN=$HOME/.nix-go/bin
    mkdir -p $GOBIN
    export PATH=$GOBIN:$PATH
    # install gotestloghelper
    go install github.com/smartcontractkit/chainlink-testing-framework/tools/gotestloghelper@latest
  '';
}

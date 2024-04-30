{
  description = "chainlink-testing-framework development shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs@{ self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; overlays = [ ]; };
      in rec {
        devShell = pkgs.callPackage ./shell.nix {
          inherit pkgs;
          scriptDir = toString ./.;  # This converts the flake's root directory to a string
        };
        formatter = pkgs.nixpkgs-fmt;
      });
}

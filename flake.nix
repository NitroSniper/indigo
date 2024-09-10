{
  # https://wiki.nixos.org/wiki/Python#With_pyproject.toml
  # This code uses pyproject.toml, if you want to use something else find it here
  # https://nix-community.github.io/pyproject.nix/introduction.html
  description = "A flake for developing python application";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    inputs:
    inputs.flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import inputs.nixpkgs { inherit system; };

        utilities = with pkgs; [
          gopls
          go
          gotools
        ];
      in
      {
        packages = { };

        devShells =
          let
            util = pkgs.mkShell { packages = utilities; };
            battery = pkgs.mkShell { packages = utilities; };
            chain = pkgs.mkShell { packages = null; };
          in
          {
            inherit battery chain util;
            default = battery;
          };
      }
    );
}

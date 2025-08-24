{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    treefmt-nix,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        inherit (nixpkgs) lib;
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
        treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
      in {
        inherit lib;

        formatter = treefmtEval.config.build.wrapper;

        checks.formatting = treefmtEval.config.build.check self;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Go
            go
            gopls
            gotools

            # Web
            deno

            # Formatting
            treefmt

            # Nix
            alejandra
            deadnix

            # Building using Docker Compose with Bake
            docker-buildx
          ];

          shellHook = ''
            GO_VERSION=$(go version | awk '{print $3}')

            printf '\n> Go version: %s\n\n' "$GO_VERSION"
          '';
        };
      }
    );
}


{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
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
        treefmtEval = treefmt-nix.lib.evalModule pkgs (import ./treefmt.nix {
          inherit pkgs;
          pkgs' = {
            node-packages = pkgs.callPackage ./pkgs/node-packages {};
          };
        });
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
            golangci-lint
            go-swag

            # Web
            deno
            nodejs_24

            # Formatting
            treefmt

            # Nix
            alejandra
            deadnix

            # Git
            # husky # Outdated, will not work with new major
            git-conventional-commits
          ];

          shellHook = ''
            export COMPOSE_PROFILES="''${COMPOSE_PROFILES:-dev}";

            printf '\n> Go version:        %s\n' "$(go version | awk '{print $3}')"
            printf '> Node version:      %s\n' "$(node --version)"
            printf '> $COMPOSE_PROFILES: %s\n\n' "''${COMPOSE_PROFILES:-N/A}"

            npx husky
          '';
        };
      }
    );
}

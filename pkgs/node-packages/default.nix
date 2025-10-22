{
  pkgs,
  system,
  nodejs_24,
}: let
  nodePackages = import ./composition.nix {
    inherit pkgs system;
    nodejs = nodejs_24;
  };
in
  nodePackages

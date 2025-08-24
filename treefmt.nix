_: {
  projectRootFile = "flake.nix";
  programs = {
    # Nix
    alejandra.enable = true;
    deadnix.enable = true;

    # Go
    gofmt.enable = true;
    goimports.enable = true;

    # Web
    deno.enable = true;
  };
}


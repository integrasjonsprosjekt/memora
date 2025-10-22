{pkgs', ...}: let
  nix.excludes = ["pkgs/**/*.nix"];
  backend.includes = ["backend/**/*.go"];
  frontend.includes = ["frontend/**/*.{js,jsx,ts,tsx,md,mdx,json,yaml,yml,css,scss,html}"];
in {
  projectRootFile = "flake.nix";

  settings.global.excludes = [
    "node_modules/**"
    ".next/**"
    ".turbo/**"
    ".vercel/**"
    "dist/**"
    "build/**"
    "coverage/**"
    "**/vendor/**"
    "go/pkg/**"
    "go/bin/**"
    "frontend/src/components/ui/**"
  ];

  programs = {
    # Nix
    alejandra = {
      inherit (nix) excludes;
      enable = true;
    };
    deadnix = {
      inherit (nix) excludes;
      enable = true;
    };

    # Go
    gofmt = {
      inherit (backend) includes;
      enable = true;
      priority = 10;
    };
    goimports = {
      inherit (backend) includes;
      enable = true;
      priority = 20;
    };
    golines = {
      inherit (backend) includes;
      enable = true;
      priority = 30;
    };

    # Web
    prettier = {
      inherit (frontend) includes;
      enable = true;
      settings =
        builtins.fromJSON (builtins.readFile ./frontend/.prettierrc)
        // {
          pluginSearchDirs = [
            "frontend"
          ];
          plugins = [
            # https://github.com/numtide/treefmt-nix/issues/112#issuecomment-1691563490
            "${pkgs'.node-packages.prettier-plugin-tailwindcss}/lib/node_modules/prettier-plugin-tailwindcss/dist/index.mjs"
          ];
        };
    };
  };
}

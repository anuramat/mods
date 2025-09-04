{
  description = "mods";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"

        "aarch64-darwin"
        "x86_64-darwin"
      ];
      imports = [
        inputs.treefmt-nix.flakeModule
      ];
      perSystem =
        {
          self',
          pkgs,
          ...
        }:
        let
          package = pkgs.buildGoModule {
            pname = "mods";
            version = "unstable";
            src = ./.;
            vendorHash = "sha256-c6uiuN48dwO7Lma2jQdRflWYyErbUmYN1vJqO7LOlU4=";

            nativeBuildInputs = [ pkgs.installShellFiles ];

            postInstall = ''
              export HOME=$(mktemp -d)
              export XDG_CONFIG_HOME=$HOME/.config
              installShellCompletion --cmd mods \
                --bash <($out/bin/mods completion bash) \
                --fish <($out/bin/mods completion fish) \
                --zsh <($out/bin/mods completion zsh)
            '';

            meta.mainProgram = "mods";
            doCheck = false;
          };
        in
        {
          packages.default = package;
          treefmt.config.programs = {
            gofumpt.enable = true;
            goimports.enable = true;
            nixfmt.enable = true;
          };
        };
    };
}

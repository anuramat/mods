{
  description = "mods";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "mods";
          version = "unstable";
          src = ./.;
          vendorHash = "sha256-c6uiuN48dwO7Lma2jQdRflWYyErbUmYN1vJqO7LOlU4=";
          meta.mainProgram = "mods";
          doCheck = false;
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/mods";
        };
      }
    );
}

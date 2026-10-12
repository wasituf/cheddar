{
  description = "Go flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";

    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      gomod2nix,
    }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ gomod2nix.overlays.default ];
      };
      ldLibs = with pkgs; [
        libX11
        libXcursor
        libXi
        libXinerama
        libXrandr
        libXrender
        libXext
        libGL
        alsa-lib
      ];
      goPackage = pkgs.buildGoApplication {
        pname = "cheddar";
        version = "0.1.0";
        src = ./.;
        modules = ./gomod2nix.toml;
        nativeBuildInputs = with pkgs; [
          pkg-config
        ];
        buildInputs = ldLibs;
      };
    in
    {
      packages.${system}.default = goPackage;
      devShells.x86_64-linux.default = pkgs.mkShell {
        packages = with pkgs; [
          git
          go
          gopls
          gotestsum
          gomod2nix.packages.${system}.default
          go-task
        ];

        buildInputs = ldLibs;

        shellHook = ''
          export NIX_LD_LIBRARY_PATH="${pkgs.lib.makeLibraryPath ldLibs}"
          export LD_LIBRARY_PATH="${pkgs.lib.makeLibraryPath ldLibs}:$LD_LIBRARY_PATH"
        '';
      };
    };
}

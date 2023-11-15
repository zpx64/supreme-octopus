{
  inputs = {
    nixpkgs.url = "nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        packages = with pkgs; [
          bun
          nodejs_21
          nodePackages_latest.typescript-language-server
        ];
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = packages;
        };
      });
}

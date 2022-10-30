{
  description = "Purple DB";

  outputs = { self, nixpkgs, flake-utils }:
    let
      goOverlay = self: super: {
        go = super.go_1_18;
      };
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; overlays = [ goOverlay ]; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs;
            [
              go
              protobuf
              docker
            ];
        };
      });
}

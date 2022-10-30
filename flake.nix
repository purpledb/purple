{
  description = "Purple DB";

  outputs = { self, nixpkgs, flake-utils }:
    let
      # Constants
      target = {
        os = "linux";
        arch = "amd64";
      };

      # The same across all packages
      vendorSha256 = "sha256-rRHPYuWOPIHE60vVxZMH+TwFNGlQf5fPmoqV+g0cUZg=";

      # Overlays
      goOverlay = self: super: {
        go = super.go_1_18;
      };

      goLinuxOverlay = self: super: {
        buildGoModuleLinux = super.buildGoModule.override {
          go = super.go // {
            CGO_ENABLED = 0;
            GOOS = target.os;
            GOARCH = target.arch;
          };
        };
      };
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; overlays = [ goOverlay goLinuxOverlay ]; };
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

        packages =
          let
            buildGo = name: pkgs.buildGoModule {
              inherit name vendorSha256;
              src = ./.;
              subPackages = [ "cmd/${name}" ];
            };

            buildDocker = { name, bin, port }:
              let
                pkg = pkgs.buildGoModuleLinux {
                  inherit name vendorSha256;
                  src = ./.;
                  subPackages = [ "cmd/${bin}" ];
                  maxLayers = 120;
                };
              in
              pkgs.dockerTools.buildLayeredImage {
                inherit name;
                tag = "latest";

                config = {
                  Entrypoint = [ "${pkg}/bin/${target.os}_${target.arch}/${bin}" ];
                  ExposedPorts."${builtins.toString port}/tcp" = { };
                };
              };
          in
          rec
          {
            http = purpleHttp;
            grpc = purpleGrpc;

            purpleHttp = buildGo "purple-http";
            purpleGrpc = buildGo "purple-grpc";

            purpleHttpDocker = buildDocker rec {
              name = "purpledb/${bin}";
              bin = "purple-http";
              port = 8080;
            };

            purpleGrpcDocker = buildDocker rec {
              name = "purpledb/${bin}";
              bin = "purple-grpc";
              port = 8081;
            };
          };
      });
}

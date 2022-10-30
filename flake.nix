{
  description = "Purple DB";

  outputs = { self, nixpkgs, flake-utils }:
    let
      # Constants
      org = "purpledb";
      version = "0.1.6";

      target = {
        os = "linux";
        arch = "amd64";
      };

      # The same across all packages
      vendorSha256 = "sha256-rRHPYuWOPIHE60vVxZMH+TwFNGlQf5fPmoqV+g0cUZg=";

      # Overlays
      overlays = [
        (self: super: rec {
          go = super.go_1_18;

          buildGoModuleLinux = super.buildGoModule.override {
            go = super.go // {
              CGO_ENABLED = 0;
              GOOS = target.os;
              GOARCH = target.arch;
            };
          };
        })
      ];
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit overlays system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs;
            [
              cmake
              go
              gotools
              protobuf
              docker
            ];
        };

        apps = {
          default = flake-utils.lib.mkApp {
            drv = pkgs.writeScriptBin "version" ''
              echo ${version}
            '';
          };
        };

        packages =
          let
            buildGo = name: pkgs.buildGoModule {
              inherit name vendorSha256;
              src = ./.;
              subPackages = [ "cmd/${name}" ];
            };

            buildDocker = { name, port, tag ? "latest" }:
              let
                pkg = pkgs.buildGoModuleLinux {
                  inherit name vendorSha256;
                  src = ./.;
                  subPackages = [ "cmd/${name}" ];
                  maxLayers = 120;
                };
              in
              pkgs.dockerTools.buildLayeredImage {
                name = "${org}/${name}";
                inherit tag;

                config = {
                  Entrypoint = [ "${pkg}/bin/${target.os}_${target.arch}/${name}" ];
                  ExposedPorts."${builtins.toString port}/tcp" = { };
                };
              };
          in
          rec
          {

            http = buildGo "purple-http";
            grpc = buildGo "purple-grpc";

            httpDocker = buildDocker rec {
              name = "purple-http";
              port = 8080;
              tag = "v${version}";
            };

            httpDockerLatest = buildDocker rec {
              name = "purple-http";
              port = 8080;
            };

            grpcDocker = buildDocker rec {
              name = "purple-grpc";
              port = 8081;
              tag = "v${version}";
            };

            grpcDockerLatest = buildDocker rec {
              name = "purple-grpc";
              port = 8081;
            };
          };
      });
}

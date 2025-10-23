{
  description = "Purple DB";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1";

  outputs =
    {
      self,
      ...
    }@inputs:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-darwin"
        "x86_64-darwin"
        "aarch64-linux"
      ];
      forEachSupportedSystem =
        f:
        inputs.nixpkgs.lib.genAttrs supportedSystems (
          system:
          f {
            pkgs = import inputs.nixpkgs {
              inherit system;
              overlays = [ self.overlays.default ];
            };
          }
        );

      # Constants
      org = "purpledb";
      version = "0.1.6";

      target = {
        os = "linux";
        arch = "amd64";
      };

      # The same across all packages
      vendorHash = "sha256-fSdq6j63GbTHAZE93gtHeMPYCLT3PdiVrG4aqWAD2QE=";
    in
    {
      devShells = forEachSupportedSystem (
        { pkgs }:
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              cmake
              go
              gotools
              protobuf
              docker
            ];
          };
        }
      );

      packages = forEachSupportedSystem (
        { pkgs }:
        {
          http = pkgs.buildGo "purple-http";
          grpc = pkgs.buildGo "purple-grpc";
        }
      );

      dockerImages = forEachSupportedSystem (
        { pkgs }:
        {
          httpDocker = pkgs.buildDocker {
            name = "purple-http";
            port = 8080;
            tag = "v${version}";
          };

          httpDockerLatest = pkgs.buildDocker {
            name = "purple-http";
            port = 8080;
          };

          grpcDocker = pkgs.buildDocker {
            name = "purple-grpc";
            port = 8081;
            tag = "v${version}";
          };

          grpcDockerLatest = pkgs.buildDocker {
            name = "purple-grpc";
            port = 8081;
          };
        }
      );

      overlays.default = final: prev: {
        buildGoModuleLinux = prev.buildGoModule.override {
          go = prev.go // {
            CGO_ENABLED = 0;
            GOOS = target.os;
            GOARCH = target.arch;
          };
        };

        buildGo =
          name:
          final.buildGoModule {
            inherit name vendorHash;
            src = ./.;
            subPackages = [ "cmd/${name}" ];
          };

        buildDocker =
          {
            name,
            port,
            tag ? "latest",
          }:
          let
            pkg = final.buildGoModuleLinux {
              inherit name vendorHash;
              src = ./.;
              subPackages = [ "cmd/${name}" ];
              maxLayers = 120;
            };
          in
          final.dockerTools.buildLayeredImage {
            name = "${org}/${name}";
            inherit tag;

            config = {
              Entrypoint = [ "${pkg}/bin/${target.os}_${target.arch}/${name}" ];
              ExposedPorts."${builtins.toString port}/tcp" = { };
            };
          };
      };
    };
}

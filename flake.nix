{
  description = "Generate docs from scripted test commands";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nix-filter.url = "github:numtide/nix-filter";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, nix-filter }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        localOverlay = import ./nix/overlay.nix;
        overlays = [ localOverlay ];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        go-src = nix-filter.lib.filter {
          root = ./.;
          include = with nix-filter.lib;[
            (nix-filter.lib.matchExt "go")
            ./go.mod
            ./go.sum
            (inDirectory "vendor")
            isDirectory
          ];
        };

        nix-src = nix-filter.lib.filter {
          root = ./.;
          include = [
            (nix-filter.lib.matchExt "nix")
          ];
        };

        buildInputs = with pkgs; [
        ];

        nativeBuildInputs = with pkgs; [
          go
        ];

        name = "tad";
        description = description;
        version = nixpkgs.lib.fileContents ./VERSION;
        module = "github.com/nhost/tad";

        tags = "integration";

        ldflags = ''
          -X main.Version=${version}
        '';

      in
      {
        checks = {
          nixpkgs-fmt = pkgs.runCommand "check-nixpkgs-fmt"
            {
              nativeBuildInputs = with pkgs;
                [
                  nixpkgs-fmt
                ];
            }
            ''
              mkdir $out
              nixpkgs-fmt --check ${nix-src}
            '';

          linters = pkgs.runCommand "linters"
            {
              nativeBuildInputs = with pkgs; [
                govulncheck
                golangci-lint
              ] ++ buildInputs ++ nativeBuildInputs;
            }
            ''
              export GOLANGCI_LINT_CACHE=$TMPDIR/.cache/golangci-lint
              export GOCACHE=$TMPDIR/.cache/go-build
              export GOMODCACHE="$TMPDIR/.cache/mod"

              mkdir $out
              cd $out
              cp -r ${go-src}/* .

              govulncheck ./...

              golangci-lint run \
                --build-tags=${tags} \
                --timeout 300s
            '';

          gotests = pkgs.runCommand "gotests"
            {
              nativeBuildInputs = with pkgs; [
                go
              ] ++ buildInputs ++ nativeBuildInputs;
            }
            ''
              export GOCACHE=$TMPDIR/.cache/go-build
              export GOMODCACHE="$TMPDIR/.cache/mod"

              mkdir $out
              cd $out
              cp -r ${go-src}/* .

              export GIN_MODE=release

              go test \
                -tags=${tags} \
                -ldflags="${ldflags}" \
                -v ./...
            '';

        };

        devShells = flake-utils.lib.flattenTree rec {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              golines
              gofumpt
              nhost-cli
              nixpkgs-fmt
              golangci-lint
              docker-client
              govulncheck
              gnumake
              richgo
            ] ++ buildInputs ++ nativeBuildInputs;
          };
        };

        packages = flake-utils.lib.flattenTree rec {
          tad = pkgs.buildGoModule {
            inherit version ldflags buildInputs nativeBuildInputs;

            pname = name;
            src = go-src;

            vendorSha256 = null;

            doCheck = false;

            subPackages = [
              "."
            ];

            meta = with pkgs.lib; {
              description = description;
              homepage = "https://github.com/nhost/tad";
              maintainers = [ "nhost" ];
              platforms = platforms.linux ++ platforms.darwin;
            };
          };

          dockerImage = pkgs.dockerTools.buildLayeredImage {
            name = name;
            tag = version;
            created = "now";
            contents = with pkgs; [
              (writeTextFile {
                name = "tmp-file";
                text = ''
                  dummy file to generate tmpdir
                '';
                destination = "/tmp/tmp-file";
              })
              cacert
            ] ++ buildInputs;
            config = {
              Env = [
                "TMPDIR=/tmp"
              ];
              Entrypoint = [
                "${self.packages.${system}.tad}/bin/tad"
              ];
            };
          };

          default = tad;

        };

      }



    );


}

final: prev: rec {
  nhost-cli = import ./nhost-cli.nix { inherit final prev; };

  go = prev.go_1_21.overrideAttrs
    (finalAttrs: previousAttrs: rec {
      version = "1.21.3";

      src = final.fetchurl {
        url = "https://go.dev/dl/go${version}.src.tar.gz";
        sha256 = "sha256-GG8rb4yLcE5paCGwmrIEGlwe4T3LwxVqE63PdZMe5Ig=";
      };

    });

  buildGoModule = prev.buildGoModule.override { go = go; };

  golangci-lint = prev.golangci-lint.override rec {
    buildGoModule = args: final.buildGoModule (args // rec {
      version = "1.53.2";
      src = final.fetchFromGitHub {
        owner = "golangci";
        repo = "golangci-lint";
        rev = "v${version}";
        sha256 = "sha256-fsK9uHPh3ltZpAlo4kDp9MtGKBYqDMxJSygvaloKbMs=";
      };
      vendorHash = "sha256-gg4E+6u1aukyXWJrRJRmhdeqUN/zw76kznSuE+99CDc=";

      meta = with final.lib; args.meta // {
        broken = false;
      };
    });
  };

  golines = final.buildGoModule rec {
    name = "golines";
    version = "0.11.0";
    src = final.fetchFromGitHub {
      owner = "dbarrosop";
      repo = "golines";
      rev = "b7e767e781863a30bc5a74610a46cc29485fb9cb";
      sha256 = "sha256-pxFgPT6J0vxuWAWXZtuR06H9GoGuXTyg7ue+LFsRzOk=";
    };
    vendorHash = "sha256-rxYuzn4ezAxaeDhxd8qdOzt+CKYIh03A9zKNdzILq18=";

    meta = with final.lib; {
      description = "A golang formatter that fixes long lines";
      homepage = "https://github.com/segmentio/golines";
      maintainers = [ "nhost" ];
      platforms = platforms.linux ++ platforms.darwin;
    };
  };

  govulncheck = final.buildGoModule rec {
    name = "govulncheck";
    version = "v1.0.1";
    src = final.fetchFromGitHub {
      owner = "golang";
      repo = "vuln";
      rev = "${version}";
      sha256 = "sha256-wAZJruDyAavfqS/+7yDaT4suE4eivCbL7g4tQoqUDlQ=";
    };
    vendorHash = "sha256-LqKDQPOXOF+A2G3e483jeV7ziUb87ePnk8VMkLAms74=";

    doCheck = false;

    meta = with final.lib; {
      description = "the database client and tools for the Go vulnerability database";
      homepage = "https://github.com/golang/vuln";
      maintainers = [ "nhost" ];
      platforms = platforms.linux ++ platforms.darwin;
    };
  };

  gqlgen = final.buildGoModule rec {
    pname = "gqlgen";
    version = "0.17.36";

    src = final.fetchFromGitHub {
      owner = "99designs";
      repo = pname;
      rev = "v${version}";
      sha256 = "sha256-ieJfnBEJEGBQJONm4YfFhbRY1pED5gcsfbcDemqAMLE=";
    };

    vendorHash = "sha256-mYCPRfcTZNc7i7mydDO9/9pXJJ2bB7HhKAnxJD7qxz0=";

    doCheck = false;

    subPackages = [ "./." ];

    meta = with final.lib; {
      description = "go generate based graphql server library";
      homepage = "https://gqlgen.com";
      license = licenses.mit;
      maintainers = [ "@nhost" ];
    };
  };

  gqlgenc = final.buildGoModule rec {
    pname = "gqlgenc";
    version = "0.14.0";

    src = final.fetchFromGitHub {
      owner = "Yamashou";
      repo = pname;
      rev = "v${version}";
      sha256 = "sha256-0KUJlz8ey0kLmHO083ZPaJYIhInlKvO/a1oZYjPGopo=";
    };

    vendorHash = "sha256-Up7Wi6z0Cbp9RHKAsjj/kd50UqcXtsS+ETRYuxRfGuA=";

    doCheck = false;

    subPackages = [ "./." ];

    meta = with final.lib; {
      description = "This is Go library for building GraphQL client with gqlgen";
      homepage = "https://github.com/Yamashou/gqlgenc";
      license = licenses.mit;
      maintainers = [ "@nhost" ];
    };
  };

}

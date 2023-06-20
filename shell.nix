with (import <nixpkgs> {});
mkShell {
  buildInputs = [
    go_1_20
    golint
    gopls
    sqlite
    flyctl
    golangci-lint
    foreman
    google-cloud-sdk
  ];
}

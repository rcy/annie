let
  unstable = import (fetchTarball https://nixos.org/channels/nixos-unstable/nixexprs.tar.xz) { };
in
{ nixpkgs ? import <nixpkgs> {} }:
with nixpkgs; mkShell {
  buildInputs = [
    air
    unstable.go_1_21
    unstable.golint
    unstable.gopls
    html2text
    sqlite
    flyctl
    unstable.golangci-lint
    google-cloud-sdk
  ];
}

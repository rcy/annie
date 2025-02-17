let
  unstable = import (fetchTarball https://nixos.org/channels/nixos-unstable/nixexprs.tar.xz) { };
in
{ nixpkgs ? import <nixpkgs> {} }:
with nixpkgs; mkShell {
  buildInputs = [
    nodejs
    air
    go
    golint
    gopls
    html2text
    flyctl
    golangci-lint
    google-cloud-sdk
    pup
  ];
}

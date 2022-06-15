with (import <nixpkgs> {});
mkShell {
  buildInputs = [
    go
    sqlite
    flyctl
  ];
}

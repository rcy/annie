with (import <nixpkgs> {});
mkShell {
  buildInputs = [
    pipenv
  ];
}

{  
  description = "";  
  
  inputs = {  
    nixpkgs.url = "nixpkgs";  
  };  
  
  outputs = { self, nixpkgs }:   
  let   
      # Systems supported  
      allSystems = [  
        "x86_64-linux" # 64-bit Intel/AMD Linux  
        "aarch64-linux" # 64-bit ARM Linux  
        "x86_64-darwin" # 64-bit Intel macOS  
        "aarch64-darwin" # 64-bit ARM macOS  
      ];  
  
      # Helper to provide system-specific attributes  
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {  
        pkgs = import nixpkgs { inherit system; };  
      });  
  in {  
    # Development env package required.  
    overlays.default = final: prev: rec {
      nodejs = prev.nodejs;
      yarn = (prev.yarn.override { inherit nodejs; });
    };
    devShells = forAllSystems ({ pkgs }: {  
        default = pkgs.mkShell {  
          # The Nix packages provided in the environment  
          hardeningDisable = ["fortify"];
          packages = with pkgs; [  
            prettierd
            air
            gofumpt
            goimports-reviser
            golines
            sqlitebrowser
            go_1_23 # Go 1.22  
            sqlitebrowser
            gotools # Go tools like goimports, godoc, and others  
            delve
            node2nix 
            gcc
            gdb
            pkg-config
            nodejs
          ];  
shellHook = ''
  export CGO_ENABLED=1
  export CC=${pkgs.gcc}/bin/gcc
  export CXX=${pkgs.gcc}/bin/g++
'';

        };  
    });  
  };  
}


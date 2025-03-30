# moddownloader-go

A simple Mod updater using the Modrinth API.

### Usage

The program can be used either interactively via the command line (CLI) or with arguments.

- `-help`: Shows help.
- `-version`: Sets the Minecraft version.
- `-loader`: Sets the mod loader (e.g., "fabric").
- `-input`: Specifies the input folder containing the mods to be updated.
- `-output`: Specifies the output folder where the downloaded mods should be saved.

Example (Windows):
```console
.\moddownloader-go.exe -version 1.21.5 -loader fabric -input mods_to_update -output out
```

### Download

The latest release can be found [here](https://github.com/Deskilling/moddownloader-go/releases).

### Building

1.  Clone the Repository:

    ```console
    git clone https://github.com/Deskilling/moddownloader-go.git
    ```

2.  Open the folder in the terminal.

   3.  Compile the program using the Go build tools:

       ```console
       go build -o .
       ```